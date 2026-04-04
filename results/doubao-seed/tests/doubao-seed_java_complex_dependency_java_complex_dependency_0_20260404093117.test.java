import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnit;
import org.mockito.junit.MockitoRule;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import org.apache.commons.lang3.SystemUtils;
import org.apache.commons.lang3.StringUtils;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import static org.junit.Assert.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;

@RunWith(PowerMockRunner.class)
@PrepareForTest({System.class, SystemUtils.class, StringUtils.class, LocalAndroidPlatforms.class})
public class LocalAndroidPlatformsTest {

    @Rule
    public MockitoRule mockitoRule = MockitoJUnit.rule();
    @Rule
    public org.junit.rules.ExpectedException thrown = org.junit.rules.ExpectedException.none();

    @Mock
    private File mockSdkDir;
    @Mock
    private File mockInvalidSdkDir;
    @Mock
    private File mockPlatformToolsDir;
    @Mock
    private File mockAdbFile;
    @Mock
    private LocalAndroidPlatforms mockValidPlatforms;
    @Mock
    private LocalAndroidPlatforms mockInvalidPlatforms;

    @Test
    public void testFindLocalJavaSdk_androidHomeEnvVarValid_returnsSdkDir() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn("/test/android/sdk");
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(eq(mockSdkDir)).thenReturn(mockValidPlatforms);
        PowerMockito.when(mockValidPlatforms.valid()).thenReturn(true);
        PowerMockito.whenNew(File.class).withArguments(eq("/test/android/sdk")).thenReturn(mockSdkDir);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(mockSdkDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_androidHomeInvalidAndroidSdkRootValid_returnsSdkDir() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn("/test/invalid/sdk");
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn("/test/valid/sdk");
        PowerMockito.whenNew(File.class).withArguments(eq("/test/invalid/sdk")).thenReturn(mockInvalidSdkDir);
        PowerMockito.whenNew(File.class).withArguments(eq("/test/valid/sdk")).thenReturn(mockSdkDir);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(eq(mockInvalidSdkDir)).thenReturn(mockInvalidPlatforms);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(eq(mockSdkDir)).thenReturn(mockValidPlatforms);
        PowerMockito.when(mockInvalidPlatforms.valid()).thenReturn(false);
        PowerMockito.when(mockValidPlatforms.valid()).thenReturn(true);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(mockSdkDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_envVarsMissingWindowsAdbFoundValid_returnsSdkDir() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.mockStatic(SystemUtils.class);
        PowerMockito.mockStatic(StringUtils.class);
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.when(SystemUtils.IS_OS_WINDOWS).thenReturn(true);
        String testPath = "C:\\test\\android\\sdk\\platform-tools";
        PowerMockito.when(System.getenv("PATH")).thenReturn(testPath);
        PowerMockito.when(StringUtils.split(eq(testPath), eq(File.pathSeparatorChar))).thenReturn(new String[]{testPath});

        PowerMockito.whenNew(File.class).withArguments(eq(testPath)).thenReturn(mockPlatformToolsDir);
        PowerMockito.whenNew(File.class).withArguments(eq(mockPlatformToolsDir), eq("adb.exe")).thenReturn(mockAdbFile);
        PowerMockito.when(mockAdbFile.exists()).thenReturn(true);
        PowerMockito.when(mockAdbFile.isFile()).thenReturn(true);
        PowerMockito.when(mockAdbFile.canExecute()).thenReturn(true);
        PowerMockito.when(mockAdbFile.getParentFile()).thenReturn(mockPlatformToolsDir);
        PowerMockito.when(mockPlatformToolsDir.getParentFile()).thenReturn(mockSdkDir);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(eq(mockSdkDir)).thenReturn(mockValidPlatforms);
        PowerMockito.when(mockValidPlatforms.valid()).thenReturn(true);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(mockSdkDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_envVarsMissingNonWindowsAdbFoundValid_returnsSdkDir() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.mockStatic(SystemUtils.class);
        PowerMockito.mockStatic(StringUtils.class);
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.when(SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        String testPath = "/test/android/sdk/platform-tools";
        PowerMockito.when(System.getenv("PATH")).thenReturn(testPath);
        PowerMockito.when(StringUtils.split(eq(testPath), eq(File.pathSeparatorChar))).thenReturn(new String[]{testPath});

        PowerMockito.whenNew(File.class).withArguments(eq(testPath)).thenReturn(mockPlatformToolsDir);
        PowerMockito.whenNew(File.class).withArguments(eq(mockPlatformToolsDir), eq("adb")).thenReturn(mockAdbFile);
        PowerMockito.when(mockAdbFile.exists()).thenReturn(true);
        PowerMockito.when(mockAdbFile.isFile()).thenReturn(true);
        PowerMockito.when(mockAdbFile.canExecute()).thenReturn(true);
        PowerMockito.when(mockAdbFile.getParentFile()).thenReturn(mockPlatformToolsDir);
        PowerMockito.when(mockPlatformToolsDir.getParentFile()).thenReturn(mockSdkDir);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(eq(mockSdkDir)).thenReturn(mockValidPlatforms);
        PowerMockito.when(mockValidPlatforms.valid()).thenReturn(true);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(mockSdkDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_noEnvVarsNoValidBinaries_throwsFileNotFoundException() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.mockStatic(SystemUtils.class);
        PowerMockito.mockStatic(StringUtils.class);
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.when(SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        String testPath = "/test/random/path";
        PowerMockito.when(System.getenv("PATH")).thenReturn(testPath);
        PowerMockito.when(StringUtils.split(eq(testPath), eq(File.pathSeparatorChar))).thenReturn(new String[]{testPath});
        PowerMockito.whenNew(File.class).withArguments(any(String.class)).thenReturn(new File(testPath));
        PowerMockito.whenNew(File.class).withArguments(any(File.class), any(String.class)).thenReturn(new File(testPath, "test"));

        thrown.expect(FileNotFoundException.class);
        thrown.expectMessage("Unable to find the Local Android Java SDK Folder");
        LocalAndroidPlatforms.findLocalJavaSdk();
    }

    @Test
    public void testFindLocalJavaSdk_firstBinaryInvalidSecondValid_returnsValidDir() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.mockStatic(SystemUtils.class);
        PowerMockito.mockStatic(StringUtils.class);
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.when(SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        String testPath = "/test/android/sdk/platform-tools";
        PowerMockito.when(System.getenv("PATH")).thenReturn(testPath);
        PowerMockito.when(StringUtils.split(eq(testPath), eq(File.pathSeparatorChar))).thenReturn(new String[]{testPath});

        File mockEmulatorFile = PowerMockito.mock(File.class);
        File mockEmulatorDir = PowerMockito.mock(File.class);
        File mockValidEmulatorHome = PowerMockito.mock(File.class);

        // Mock invalid adb
        PowerMockito.whenNew(File.class).withArguments(eq(mockPlatformToolsDir), eq("adb")).thenReturn(mockAdbFile);
        PowerMockito.when(mockAdbFile.exists()).thenReturn(true);
        PowerMockito.when(mockAdbFile.isFile()).thenReturn(true);
        PowerMockito.when(mockAdbFile.canExecute()).thenReturn(true);
        PowerMockito.when(mockAdbFile.getParentFile()).thenReturn(mockPlatformToolsDir);
        PowerMockito.when(mockPlatformToolsDir.getParentFile()).thenReturn(mockInvalidSdkDir);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(eq(mockInvalidSdkDir)).thenReturn(mockInvalidPlatforms);
        PowerMockito.when(mockInvalidPlatforms.valid()).thenReturn(false);

        // Mock valid emulator
        PowerMockito.whenNew(File.class).withArguments(eq(mockPlatformToolsDir), eq("emulator")).thenReturn(mockEmulatorFile);
        PowerMockito.when(mockEmulatorFile.exists()).thenReturn(true);
        PowerMockito.when(mockEmulatorFile.isFile()).thenReturn(true);
        PowerMockito.when(mockEmulatorFile.canExecute()).thenReturn(true);
        PowerMockito.when(mockEmulatorFile.getParentFile()).thenReturn(mockEmulatorDir);
        PowerMockito.when(mockEmulatorDir.getParentFile()).thenReturn(mockValidEmulatorHome);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(eq(mockValidEmulatorHome)).thenReturn(mockValidPlatforms);
        PowerMockito.when(mockValidPlatforms.valid()).thenReturn(true);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(mockValidEmulatorHome, result);
    }
}