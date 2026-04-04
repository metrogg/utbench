import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mockito;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.util.logging.Logger;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;
import static org.mockito.Matchers.any;
import static org.mockito.Matchers.anyString;
import static org.mockito.Mockito.when;

@RunWith(PowerMockRunner.class)
@PrepareForTest({LocalAndroidPlatforms.class, SystemUtils.class, StringUtils.class, System.class, File.class})
public class LocalAndroidPlatformsTest {

    private LocalAndroidPlatforms mockPlatform;
    private File mockHomeDir;

    @Before
    public void setUp() throws Exception {
        // Mock static dependencies
        PowerMockito.mockStatic(SystemUtils.class);
        PowerMockito.mockStatic(StringUtils.class);
        PowerMockito.mockStatic(System.class);

        // Mock the instance of LocalAndroidPlatforms that the static method creates
        mockPlatform = Mockito.mock(LocalAndroidPlatforms.class);
        mockHomeDir = Mockito.mock(File.class);
        
        // Mock the constructor to return our mock instance
        PowerMockito.whenNew(LocalAndroidPlatforms.class).withAnyArguments().thenReturn(mockPlatform);
        
        // Mock the valid() method to return true by default for success cases
        when(mockPlatform.valid()).thenReturn(true);
        
        // Handle the potential field access 'platforms.valid' if the source code inconsistency exists
        try {
            PowerMockito.field(LocalAndroidPlatforms.class, "valid").set(mockPlatform, true);
        } catch (Exception e) {
            // Field might not exist or be accessible, ignore if only method is used
        }
    }

    @Test
    public void testFindLocalJavaSdk_SuccessViaAndroidHome() throws Exception {
        // Arrange
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn("/usr/lib/android-sdk");
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        
        // Act
        File result = LocalAndroidPlatforms.findLocalJavaSdk();
        
        // Assert
        assertEquals(new File("/usr/lib/android-sdk"), result);
    }

    @Test
    public void testFindLocalJavaSdk_SuccessViaAndroidSdkRoot() throws Exception {
        // Arrange
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn("/usr/lib/android-sdk-root");
");

        // Act
        File result = LocalAndroidPlatforms.findLocalJavaSdk();

        // Assert
        assertEquals(new File("/usr/lib/android-sdk-root"), result);
    }

    @Test
    public void testFindLocalJavaSdk_FallbackToSdkRootIfHomeInvalid() throws Exception {
        // Arrange
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn("/invalid/home");
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn("/valid/root");

        // First call (for ANDROID_HOME) is invalid, second (for SDK_ROOT) is valid
        when(mockPlatform.valid()).thenReturn(false, true);

        // Act
        File result = LocalAndroidPlatforms.findLocalJavaSdk();

        // Assert
        assertEquals(new File("/valid/root"), result);
    }

    @Test
    public void testFindLocalJavaSdk_SuccessViaPathUnix() throws Exception {
        // Arrange
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.when(SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        
        String pathEnv = "/usr/local/bin:/usr/bin";
        PowerMockito.when(System.getenv("PATH")).thenReturn(pathEnv);
        PowerMockito.when(StringUtils.split(pathEnv, File.pathSeparatorChar))
                .thenReturn(new String[]{"/usr/local/bin", "/usr/bin"});

        // Mock File interactions for the search
        File mockPathDir = Mockito.mock(File.class);
        File mockBinFile = Mockito.mock(File.class);
        File mockBinParent = Mockito.mock(File.class);

        // Intercept new File(pathPart)
        PowerMockito.whenNew(File.class).withArguments("/usr/local/bin").thenReturn(mockPathDir);
        // Intercept new File(dir, "adb")
        PowerMockito.whenNew(File.class).withArguments(mockPathDir, "adb").thenReturn(mockBinFile);

        when(mockBinFile.exists()).thenReturn(true);
        when(mockBinFile.isFile()).()).thenReturn(true);
        when(mockBinFile.canExecute()).thenReturn(true);
        when(mockBinFile.getParentFile()).thenReturn(mockBinParent);
        when(mockBinParent.getParentFile()).thenReturn(mockHomeDir);

        // Act
        File result = LocalAndroidPlatforms.findLocalJavaSdk();

        // Assert
        assertEquals(mockHomeDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_SuccessViaPathWindows() throws Exception {
        // Arrange
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.when(SystemUtils.IS_OS_WINDOWS).thenReturn(true);

        String pathEnv = "C:\\Android\\sdk\\platform-tools";
        PowerMockito.when(System.getenv("PATH")).thenReturn(pathEnv);
        PowerMockito.when(StringUtils.split(pathEnv, File.pathSeparatorChar))
                .thenReturn(new String[]{"C:\\Android\\sdk\\platform-tools"});

        // Mock File interactions
        File mockPathDir = Mockito.mock(File.class);
        File mockBinFile = Mockito.mock(File.class);
        File mockBinParent = Mockito.mock(File.class);

        PowerMockito.whenNew(File.class).withArguments("C:\\Android\\sdk\\platform-tools").thenReturn(mockPathDir);
        PowerMockito.whenNew(File.class).withArguments(mockPathDir, "adb.exe").thenReturn(mockBinFile);

        when(mockBinFile.exists()).thenReturn(true);
        when(mockBinFile.isFile()).thenReturn(true);
        when(mockBinFile.canExecute()).thenReturn(true);
        when(mockBinFile.getParentFile()).thenReturn(mockBinParent);
        when(mockBinParent.getParentFile()).thenReturn(mockHomeDir);

        // Act
        File result = LocalAndroidPlatforms.findLocalJavaSdk();

        // Assert
        assertEquals(mockHomeDir, result);
    }

    @Test
    public void testFindLocalJavaSdk_PathBinaryFoundButHomeInvalid() throws Exception {
        // Arrange
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.when(SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        PowerMockito.when(System.getenv("PATH")).thenReturn("/usr/bin");
        PowerMockito.when(StringUtils.split(anyString(), any(char.class))).thenReturn(new String[]{"/usr/bin"});

        // Mock File interactions
        File mockPathDir = Mockito.mock(File.class);
        File mockBinFile = Mockito.mock(File.class);
        File mockBinParent = Mockito.mock(File.class);

        PowerMockito.whenNew(File.class).withArguments("/usr/bin").thenReturn(mockPathDir);
        PowerMockito.whenNew(File.class).withArguments(mockPathDir, "adb").thenReturn(mockBinFile);

        when(mockBinFile.exists()).thenReturn(true);
        when(mockBinFile.isFile()).thenReturn(true);
        when(mockBinFile.canExecute()).thenReturn(true);
        when(mockBinFile.getParentFile()).thenReturn(mockBinParent);
        when(mockBinParent.getParentFile()).thenReturn(mockHomeDir);

        // The found home directory is invalid
        when(mockPlatform.valid()).thenReturn(false);

        // Act & Assert
        try {
            LocalAndroidPlatforms.findLocalJavaSdk();
            fail("Expected FileNotFoundException");
        } catch (FileNotFoundException e) {
            // Expected
        }
    }

    @Test(expected = FileNotFoundException.class)
    public void testFindLocalJavaSdk_NotFoundThrowsException() throws Exception {
        // Arrange
        PowerMockito.when(System.getenv("ANDROID_HOME")).thenReturn(null);
        PowerMockito.when(System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
        PowerMockito.when(SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        PowerMockito.when(System.getenv("PATH")).thenReturn("/usr/bin");
        PowerMockito.when(StringUtils.split(anyString(), any(char.class))).thenReturn(new String[]{"/usr/bin"});

        // Mock File interactions to simulate binary not found
        File mockPathDir = Mockito.mock(File.class);
        File mockBinFile = Mockito.mock(File.class);

        PowerMockito.whenNew(File.class).withArguments("/usr/bin").thenReturn(mockPathDir);
        PowerMockito.whenNew(File.class).withArguments(mockPathDir, "adb").thenReturn(mockBinFile);

        when(mockBinFile.exists()).thenReturn(false); // Not found

        // Act
        LocalAndroidPlatforms.findLocalJavaSdk();
    }
}