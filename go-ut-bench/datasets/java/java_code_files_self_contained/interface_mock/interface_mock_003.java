// Converted Java method
import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

class BudgetVersionManager {
    /**
     * Manages and processes budget version data with advanced validation and analysis.
     * This class provides methods to validate, filter, and analyze budget versions.
     */
    
    /**
     * Validates a budget version based on business rules.
     * @param version The budget version to validate
     * @return List of validation error messages (empty if valid)
     */
    public List<String> validateBudgetVersion(VersaoOrcamento version) {
        List<String> errors = new ArrayList<>();
        
        if (version == null) {
            errors.add("Budget version cannot be null");
            return errors;
        }
        
        if (version.getIdVersaoOrcamento() == null) {
            errors.add("Version ID cannot be null");
        }
        
        if (version.getDescricaoMensagem() == null || version.getDescricaoMensagem().trim().isEmpty()) {
            errors.add("Description message is required");
        }
        
        if (version.getIdtStatus() == null || !(version.getIdtStatus().equals("A") || version.getIdtStatus().equals("R") || version.getIdtStatus().equals("P"))) {
            errors.add("Status must be A (Approved), R (Rejected), or P (Pending)");
        }
        
        if (version.getDataVigencia() == null) {
            errors.add("Effective date is required");
        }
        
        return errors;
    }
    
    /**
     * Filters budget versions by status and date range
     * @param versions List of budget versions to filter
     * @param status Status to filter by (A/R/P)
     * @param fromDate Start date (inclusive)
     * @param toDate End date (inclusive)
     * @return Filtered list of budget versions
     */
    public List<VersaoOrcamento> filterVersions(List<VersaoOrcamento> versions, String status, Timestamp fromDate, Timestamp toDate) {
        if (versions == null) {
            return new ArrayList<>();
        }
        
        return versions.stream()
            .filter(v -> v != null)
            .filter(v -> status == null || status.equals(v.getIdtStatus()))
            .filter(v -> fromDate == null || (v.getDataVigencia() != null && !v.getDataVigencia().before(fromDate)))
            .filter(v -> toDate == null || (v.getDataVigencia() != null && !v.getDataVigencia().after(toDate)))
            .collect(Collectors.toList());
    }
    
    /**
     * Analyzes budget version statistics
     * @param versions List of budget versions to analyze
     * @return Analysis result containing counts by status and date ranges
     */
    public BudgetVersionAnalysis analyzeVersions(List<VersaoOrcamento> versions) {
        BudgetVersionAnalysis analysis = new BudgetVersionAnalysis();
        
        if (versions == null || versions.isEmpty()) {
            return analysis;
        }
        
        for (VersaoOrcamento version : versions) {
            if (version == null) continue;
            
            analysis.totalVersions++;
            
            if ("A".equals(version.getIdtStatus())) {
                analysis.approvedCount++;
            } else if ("R".equals(version.getIdtStatus())) {
                analysis.rejectedCount++;
            } else if ("P".equals(version.getIdtStatus())) {
                analysis.pendingCount++;
            }
            
            if (version.getDataVigencia() != null) {
                analysis.earliestDate = analysis.earliestDate == null ? version.getDataVigencia() : 
                    (version.getDataVigencia().before(analysis.earliestDate) ? version.getDataVigencia() : analysis.earliestDate);
                analysis.latestDate = analysis.latestDate == null ? version.getDataVigencia() : 
                    (version.getDataVigencia().after(analysis.latestDate) ? version.getDataVigencia() : analysis.latestDate);
            }
        }
        
        return analysis;
    }
    
    public static class BudgetVersionAnalysis {
        public int totalVersions = 0;
        public int approvedCount = 0;
        public int rejectedCount = 0;
        public int pendingCount = 0;
        public Timestamp earliestDate = null;
        public Timestamp latestDate = null;
        
        @Override
        public String toString() {
            return String.format(
                "Total: %d, Approved: %d, Rejected: %d, Pending: %d, Earliest: %s, Latest: %s",
                totalVersions, approvedCount, rejectedCount, pendingCount, 
                earliestDate != null ? earliestDate.toString() : "null",
                latestDate != null ? latestDate.toString() : "null"
            );
        }
    }
}

class VersaoOrcamento {
    private VersaoOrcamentoId idVersaoOrcamento;
    private String descricaoMensagem;
    private String idtStatus;
    private Timestamp dataEnvio;
    private Timestamp dataAprovacao;
    private Timestamp dataRejeicao;
    private Timestamp dataVigencia;
    private String codigoUsuario;
    private Timestamp dataLastrec;
    private Long numeroLeiorcamento;
    private Long numeroAnolei;
    private Long numeroLeildo;
    private Long numeroAnoleildo;
    private String idtTipolei;
    
    // Getters and setters
    public VersaoOrcamentoId getIdVersaoOrcamento() { return idVersaoOrcamento; }
    public void setIdVersaoOrcamento(VersaoOrcamentoId idVersaoOrcamento) { this.idVersaoOrcamento = idVersaoOrcamento; }
    public String getDescricaoMensagem() { return descricaoMensagem; }
    public void setDescricaoMensagem(String descricaoMensagem) { this.descricaoMensagem = descricaoMensagem; }
    public String getIdtStatus() { return idtStatus; }
    public void setIdtStatus(String idtStatus) { this.idtStatus = idtStatus; }
    public Timestamp getDataEnvio() { return dataEnvio; }
    public void setDataEnvio(Timestamp dataEnvio) { this.dataEnvio = dataEnvio; }
    public Timestamp getDataAprovacao() { return dataAprovacao; }
    public void setDataAprovacao(Timestamp dataAprovacao) { this.dataAprovacao = dataAprovacao; }
    public Timestamp getDataRejeicao() { return dataRejeicao; }
    public void setDataRejeicao(Timestamp dataRejeicao) { this.dataRejeicao = dataRejeicao; }
    public Timestamp getDataVigencia() { return dataVigencia; }
    public void setDataVigencia(Timestamp dataVigencia) { this.dataVigencia = dataVigencia; }
    public String getCodigoUsuario() { return codigoUsuario; }
    public void setCodigoUsuario(String codigoUsuario) { this.codigoUsuario = codigoUsuario; }
    public Timestamp getDataLastrec() { return dataLastrec; }
    public void setDataLastrec(Timestamp dataLastrec) { this.dataLastrec = dataLastrec; }
    public Long getNumeroLeiorcamento() { return numeroLeiorcamento; }
    public void setNumeroLeiorcamento(Long numeroLeiorcamento) { this.numeroLeiorcamento = numeroLeiorcamento; }
    public Long getNumeroAnolei() { return numeroAnolei; }
    public void setNumeroAnolei(Long numeroAnolei) { this.numeroAnolei = numeroAnolei; }
    public Long getNumeroLeildo() { return numeroLeildo; }
    public void setNumeroLeildo(Long numeroLeildo) { this.numeroLeildo = numeroLeildo; }
    public Long getNumeroAnoleildo() { return numeroAnoleildo; }
    public void setNumeroAnoleildo(Long numeroAnoleildo) { this.numeroAnoleildo = numeroAnoleildo; }
    public String getIdtTipolei() { return idtTipolei; }
    public void setIdtTipolei(String idtTipolei) { this.idtTipolei = idtTipolei; }
}

class VersaoOrcamentoId {
    // Simplified ID class for demonstration
    private Long id;
    
    public VersaoOrcamentoId(Long id) { this.id = id; }
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
}
