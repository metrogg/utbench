import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import static org.junit.Assert.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyChar;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;

import org.apache.commons.lang3.SystemUtils;
import org.apache.commons.lang3.StringUtils;

@RunWith(PowerMockRunner.class)
@PrepareForTest({ LocalAndroidPlatforms.class, System.class, SystemUtils.class, StringUtils.class })
public class LocalAndroidPlatformsTest {

    @Test
    public void testFindLocalJavaSdk_AndroidHomeEnvValid_Success() throws Exception {
        PowerMockito.mockStatic(System.class);
        File mockSdkDir = mock(File.class);
        when(System.getenv("ANDROID_HOME")).thenReturn("/test/android/home");
        when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);

        LocalAndroidPlatforms mockPlatforms = mock(LocalAndroidPlatforms.class);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(mockSdkDir).thenReturn(mockPlatforms);
        when(mockPlatforms.valid()).thenReturn(true);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(mockSdkDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_AndroidHomeInvalidAndroidSdkRootValid_Success() throws Exception {
        PowerMockito.mockStatic(System.class);
        File invalidDir = mock(File.class);
        File validDir = mock(File.class);
        when(System.getenv("ANDROID_HOME")).thenReturn("/invalid/path");
        when(System.getenv("ANDROID_SDK_ROOT")).thenReturn("/valid/android/sdk");

        LocalAndroidPlatforms invalidPlatforms = mock(LocalAndroidPlatforms.class);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(invalidDir).thenReturn(invalidPlatforms);
        when(invalidPlatforms.valid()).thenReturn(false);

        LocalAndroidPlatforms validPlatforms = mock(LocalAndroidPlatforms.class);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(validDir).thenReturn(validPlatforms);
        when(validPlatforms.valid()).thenReturn(true);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(validDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_EnvInvalidAdbInPathValid_Success() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.mockStatic(SystemUtils.class);
        PowerMockito.mockStatic(StringUtils.class);

        when(System.getenv("ANDROID_HOME")).thenReturn(null);
        when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        when(SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        String testPath = "/android/sdk/platform-tools" + File.pathSeparator + "/other/path";
        when(System.getenv("PATH")).thenReturn(testPath);
        when(StringUtils.split(testPath, File.pathSeparatorChar)).thenReturn(new String[]{"/android/sdk/platform-tools", "/other/path"});

        File pathDir = mock(File.class);
        File adbBin = mock(File.class);
        File platformToolsDir = mock(File.class);
        File sdkHomeDir = mock(File.class);

        PowerMockito.whenNew(File.class).withArguments("/android/sdk/platform-tools").thenReturn(pathDir);
        PowerMockito.whenNew(File.class).withArguments(pathDir, "adb").thenReturn(adbBin);
        when(adbBin.exists()).thenReturn(true);
        when(adbBin.isFile()).thenReturn(true);
        when(adbBin.canExecute()).thenReturn(true);
        when(adbBin.getParentFile()).thenReturn(platformToolsDir);
        when(platformToolsDir.getParentFile()).thenReturn(sdkHomeDir);

        LocalAndroidPlatforms validPlatforms = mock(LocalAndroidPlatforms.class);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(sdkHomeDir).thenReturn(validPlatforms);
        when(validPlatforms.valid()).thenReturn(true);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(sdkHomeDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_WindowsOsAdbExeFound_Success() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.mockStatic(SystemUtils.class);
        PowerMockito.mockStatic(StringUtils.class);

        when(System.getenv("ANDROID_HOME")).thenReturn(null);
        when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        when(SystemUtils.IS_OS_WINDOWS).thenReturn(true);
        String testPath = "C:\\android\\sdk\\platform-tools" + File.pathSeparator + "C:\\other";
        when(System.getenv("PATH")).thenReturn(testPath);
        when(StringUtils.split(testPath, File.pathSeparatorChar)).thenReturn(new String[]{"C:\\android\\sdk\\platform-tools", "C:\\other"});

        File pathDir = mock(File.class);
        File adbExe = mock(File.class);
        File platformToolsDir = mock(File.class);
        File sdkHomeDir = mock(File.class);

        PowerMockito.whenNew(File.class).withArguments("C:\\android\\sdk\\platform-tools").thenReturn(pathDir);
        PowerMockito.whenNew(File.class).withArguments(pathDir, "adb.exe").thenReturn(adbExe);
        when(adbExe.exists()).thenReturn(true);
        when(adbExe.isFile()).thenReturn(true);
        when(adbExe.canExecute()).thenReturn(true);
        when(adbExe.getParentFile()).thenReturn(platformToolsDir);
        when(platformToolsDir.getParentFile()).thenReturn(sdkHomeDir);

        LocalAndroidPlatforms validPlatforms = mock(LocalAndroidPlatforms.class);
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withArguments(sdkHomeDir).thenReturn(validPlatforms);
        when(validPlatforms.valid()).thenReturn(true);

        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        assertEquals(sdkHomeDir, result);
    }

    @Test(expected = FileNotFoundException.class)
    public void testFindLocalJavaSdk_NoValidSdkFound_ThrowsFileNotFoundException() throws Exception {
        PowerMockito.mockStatic(System.class);
        PowerMockito.mockStatic(SystemUtils.class);
        PowerMockito.mockStatic(StringUtils.class);

        when(System.getenv("ANDROID_HOME")).thenReturn(null);
        when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        when(SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        when(System.getenv("PATH")).thenReturn("/test/path1" + File.pathSeparator + "/test/path2");
        when(StringUtils.split(anyString(), anyChar())).thenReturn(new String[]{"/test/path1", "/test/path2"});

        File mockPathDir = mock(File.class);
        PowerMockito.whenNew(File.class).withArguments(anyString()).thenReturn(mockPathDir);
        File mockBin = mock(File.class);
        PowerMockito.whenNew(File.class).withArguments(any(File.class), anyString()).thenReturn(mockBin);
        when(mockBin.exists()).thenReturn(false);

        LocalAndroidPlatforms.findLocalJavaSdk();
    }
}