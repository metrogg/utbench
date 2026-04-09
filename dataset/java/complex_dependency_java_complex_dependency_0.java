class LocalAndroidPlatforms {

    public static File findLocalJavaSdk() throws IOException {
        StringBuilder err = new StringBuilder();
        err.append("Unable to find the Local Android Java SDK Folder.");

        // Check Environment Variables First
        String envKeys[] = { "ANDROID_HOME", "ANDROID_SDK_ROOT" };
        for (String envKey : envKeys) {
            File sdkHome = getEnvironmentVariableDir(err, envKey);
            if (sdkHome == null) {
                continue; // skip, not found on that key
            }
            LocalAndroidPlatforms platforms = new LocalAndroidPlatforms(sdkHome);
            if (platforms.valid()) {
                return sdkHome;
            }
        }

        // Check Path for possible android.exe (or similar)
        List<String> searchBins = new ArrayList<String>();
        if (SystemUtils.IS_OS_WINDOWS) {
            searchBins.add("adb.exe");
            searchBins.add("emulator.exe");
            searchBins.add("android.exe");
        } else {
            searchBins.add("adb");
            searchBins.add("emulator");
            searchBins.add("android");
        }

        String pathParts[] = StringUtils.split(System.getenv("PATH"), File.pathSeparatorChar);
        for (String searchBin : searchBins) {
            err.append("\nSearched PATH for ").append(searchBin);
            for (String pathPart : pathParts) {
                File pathDir = new File(pathPart);
                LOG.fine("Searching Path: " + pathDir);
                File bin = new File(pathDir, searchBin);
                if (bin.exists() && bin.isFile() && bin.canExecute()) {
                    File homeDir = bin.getParentFile().getParentFile();
                    LOG.fine("Possible Home Dir: " + homeDir);
                    LocalAndroidPlatforms platforms = new LocalAndroidPlatforms(homeDir);
                    if (platforms.valid) {
                        return homeDir;
                    }
                }
            }
            err.append(", not found.");
        }

        throw new FileNotFoundException(err.toString());
    }

    public  LocalAndroidPlatforms(File dir);

}

class LocalAndroidPlatformsTest {

    @Test
    public void testFindLocalJavaSdk() {
