#ifndef TENSORFLOW_CC_SAVED_MODEL_FINGERPRINTING_H_
#define TENSORFLOW_CC_SAVED_MODEL_FINGERPRINTING_H_
#include <string>
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "tensorflow/core/protobuf/fingerprint.pb.h"
namespace tensorflow::saved_model::fingerprinting {
absl::StatusOr<FingerprintDef> CreateFingerprintDef(
    absl::string_view export_dir);
absl::StatusOr<FingerprintDef> ReadSavedModelFingerprint(
    absl::string_view export_dir);
std::string Singleprint(uint64_t graph_def_program_hash,
                        uint64_t signature_def_hash,
                        uint64_t saved_object_graph_hash,
                        uint64_t checkpoint_hash);
std::string Singleprint(const FingerprintDef& fingerprint);
absl::StatusOr<std::string> Singleprint(absl::string_view export_dir);
}  
#endif  
#include "tensorflow/cc/saved_model/fingerprinting.h"
#include <cstdint>
#include <string>
#include "absl/container/btree_map.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/strings/strip.h"
#include "tensorflow/cc/saved_model/constants.h"
#include "tensorflow/core/framework/versions.pb.h"
#include "tensorflow/core/graph/regularization/simple_delete.h"
#include "tensorflow/core/graph/regularization/util.h"
#include "tensorflow/core/platform/env.h"
#include "tensorflow/core/platform/file_system_helper.h"
#include "tensorflow/core/platform/fingerprint.h"
#include "tensorflow/core/platform/path.h"
#include "tensorflow/core/platform/protobuf.h"  
#include "tensorflow/core/protobuf/fingerprint.pb.h"
#include "tensorflow/core/protobuf/meta_graph.pb.h"
#include "tensorflow/core/protobuf/saved_model.pb.h"
#include "tensorflow/core/protobuf/saved_object_graph.pb.h"
#include "tensorflow/core/util/tensor_bundle/naming.h"
#if !defined(PLATFORM_WINDOWS) && !defined(__APPLE__)
#include "tensorflow/cc/saved_model/fingerprinting_utils.h"
#include "tensorflow/tools/proto_splitter/cc/util.h"
#endif
#include "tsl/platform/errors.h"
#include "tsl/platform/statusor.h"
namespace tensorflow::saved_model::fingerprinting {
namespace {
using ::tensorflow::protobuf::Map;
using ::tensorflow::protobuf::io::CodedOutputStream;
using ::tensorflow::protobuf::io::StringOutputStream;
uint64_t HashCheckpointIndexFile(absl::string_view model_dir) {
  std::string meta_filename = MetaFilename(io::JoinPath(
      model_dir, kSavedModelVariablesDirectory, kSavedModelVariablesFilename));
  std::string data;
  absl::Status read_status =
      ReadFileToString(Env::Default(), meta_filename, &data);
  if (read_status.ok()) {
    return tensorflow::Fingerprint64(data);
  } else {
    LOG(WARNING) << "Failed to read checkpoint file: " << read_status;
    return 0;
  }
}
uint64_t HashSavedModel(const SavedModel& saved_model) {
  std::string saved_model_serialized;
  {
    StringOutputStream stream(&saved_model_serialized);
    CodedOutputStream output(&stream);
    output.SetSerializationDeterministic(true);
    saved_model.SerializeToCodedStream(&output);
  }
  return tensorflow::Fingerprint64(saved_model_serialized);
}
uint64_t RegularizeAndHashSignatureDefs(
    const Map<std::string, SignatureDef>& signature_def_map) {
  absl::btree_map<std::string, SignatureDef> sorted_signature_defs;
  sorted_signature_defs.insert(signature_def_map.begin(),
                               signature_def_map.end());
  uint64_t result_hash = 0;
  for (const auto& item : sorted_signature_defs) {
    result_hash =
        FingerprintCat64(result_hash, tensorflow::Fingerprint64(item.first));
    std::string signature_def_serialized;
    {
      StringOutputStream stream(&signature_def_serialized);
      CodedOutputStream output(&stream);
      output.SetSerializationDeterministic(true);
      item.second.SerializeToCodedStream(&output);
    }
    result_hash = FingerprintCat64(
        result_hash, tensorflow::Fingerprint64(signature_def_serialized));
  }
  return result_hash;
}
absl::StatusOr<uint64_t> RegularizeAndHashSavedObjectGraph(
    const SavedObjectGraph& object_graph_def) {
  absl::btree_map<int64_t, std::string> uid_to_function_names;
  for (const auto& [name, concrete_function] :
       object_graph_def.concrete_functions()) {
    TF_ASSIGN_OR_RETURN(int64_t uid, graph_regularization::GetSuffixUID(name));
    uid_to_function_names.insert({uid, name});
  }
  uint64_t result_hash = 0;
  for (const auto& [uid, function_name] : uid_to_function_names) {
    result_hash = FingerprintCat64(result_hash,
                                   tensorflow::Fingerprint64(absl::StripSuffix(
                                       function_name, std::to_string(uid))));
    std::string concrete_function_serialized;
    {
      StringOutputStream stream(&concrete_function_serialized);
      CodedOutputStream output(&stream);
      output.SetSerializationDeterministic(true);
      object_graph_def.concrete_functions()
          .at(function_name)
          .SerializeToCodedStream(&output);
    }
    result_hash = FingerprintCat64(
        result_hash, tensorflow::Fingerprint64(concrete_function_serialized));
  }
  return result_hash;
}
absl::StatusOr<FingerprintDef> CreateFingerprintDefPb(
    absl::string_view export_dir, std::string pb_file) {
  const int kFingerprintProducer = 1;
  SavedModel saved_model;
  TF_RETURN_IF_ERROR(ReadBinaryProto(Env::Default(), pb_file, &saved_model));
  FingerprintDef fingerprint_def;
  MetaGraphDef* metagraph = saved_model.mutable_meta_graphs(0);
  fingerprint_def.set_saved_model_checksum(HashSavedModel(saved_model));
  graph_regularization::SimpleDelete(*metagraph->mutable_graph_def());
  fingerprint_def.set_graph_def_program_hash(
      graph_regularization::ComputeHash(metagraph->graph_def()));
  fingerprint_def.set_signature_def_hash(
      RegularizeAndHashSignatureDefs(metagraph->signature_def()));
  TF_ASSIGN_OR_RETURN(
      uint64_t object_graph_hash,
      RegularizeAndHashSavedObjectGraph(metagraph->object_graph_def()));
  fingerprint_def.set_saved_object_graph_hash(object_graph_hash);
  fingerprint_def.set_checkpoint_hash(HashCheckpointIndexFile(export_dir));
  VersionDef* version = fingerprint_def.mutable_version();
  version->set_producer(kFingerprintProducer);
  return fingerprint_def;
}
}  
absl::StatusOr<FingerprintDef> CreateFingerprintDef(
    absl::string_view export_dir) {
  std::string prefix = io::JoinPath(export_dir, kSavedModelFilenamePrefix);
#if !defined(PLATFORM_WINDOWS) && !defined(__APPLE__)
  TF_ASSIGN_OR_RETURN(bool only_contains_pb,
                      tools::proto_splitter::OnlyContainsPb(prefix));
  if (only_contains_pb) {
    return CreateFingerprintDefPb(export_dir, absl::StrCat(prefix, ".pb"));
  }
  return CreateFingerprintDefCpb(export_dir, absl::StrCat(prefix, ".cpb"));
#else
  return CreateFingerprintDefPb(export_dir, absl::StrCat(prefix, ".pb"));
#endif
}
absl::StatusOr<FingerprintDef> ReadSavedModelFingerprint(
    absl::string_view export_dir) {
  const std::string fingerprint_pb_path =
      io::JoinPath(export_dir, kFingerprintFilenamePb);
  TF_RETURN_IF_ERROR(Env::Default()->FileExists(fingerprint_pb_path));
  FingerprintDef fingerprint_proto;
  absl::Status result =
      ReadBinaryProto(Env::Default(), fingerprint_pb_path, &fingerprint_proto);
  if (!result.ok()) return result;
  return fingerprint_proto;
}
std::string Singleprint(uint64_t graph_def_program_hash,
                        uint64_t signature_def_hash,
                        uint64_t saved_object_graph_hash,
                        uint64_t checkpoint_hash) {
  return std::to_string(graph_def_program_hash) + "/" +
         std::to_string(signature_def_hash) + "/" +
         std::to_string(saved_object_graph_hash) + "/" +
         std::to_string(checkpoint_hash);
}
std::string Singleprint(const FingerprintDef& fingerprint) {
  return Singleprint(
      fingerprint.graph_def_program_hash(), fingerprint.signature_def_hash(),
      fingerprint.saved_object_graph_hash(), fingerprint.checkpoint_hash());
}
absl::StatusOr<std::string> Singleprint(absl::string_view export_dir) {
  TF_ASSIGN_OR_RETURN(const FingerprintDef fingerprint_def,
                      ReadSavedModelFingerprint(export_dir));
  return Singleprint(fingerprint_def);
}
}  