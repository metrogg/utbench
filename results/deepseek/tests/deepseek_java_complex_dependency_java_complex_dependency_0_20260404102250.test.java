import org.junit.Test;
import org.junit.Before;
import org.junit.After;
import org.junit.Rule;
import org.junit.rules.ExpectedException;
import org.junit.rules.TemporaryFolder;
import static org.junit.Assert.*;
import static org.mockito.Mockito.*;
import org.mockito.Mockito;
import org.mockito.MockedStatic;
import org.mockito.MockitoAnnotations;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class LocalAndroidPlatformsTest {
    
    @Rule
    public TemporaryFolder tempFolder = new TemporaryFolder();
    
    @Rule
    public ExpectedException thrown = ExpectedException.none();
    
    private MockedStatic<SystemUtils> mockedSystemUtils;
    private MockedStatic<StringUtils> mockedStringUtils;
    
    @Before
    public void setUp() {
        MockitoAnnotations.initMocks(this);
    }
    
    @After
    public void tearDown() {
        if (mockedSystemUtils != null) {
            mockedSystemUtils.close();
        }
        if (mockedStringUtils != null) {
            mockedStringUtils.close();
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_AndroidHomeValid() throws IOException {
        // Mock environment variable
        File mockSdkDir = tempFolder.newFolder("android-sdk");
        File platformsDir = new File(mockSdkDir, "platforms");
        platformsDir.mkdirs();
        
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(mockSdkDir.getAbsolutePath());
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
            
            // Mock LocalAndroidPlatforms constructor and valid() method
            LocalAndroidPlatforms mockPlatforms = mock(LocalAndroidPlatforms.class);
            when(mockPlatforms.valid()).thenReturn(true);
            
            try (MockedStatic<LocalAndroidPlatforms> mockedLocalAndroidPlatforms = 
                 Mockito.mockStatic(LocalAndroidPlatforms.class)) {
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(any(File.class)))
                    .thenReturn(mockPlatforms);
                
                File result = LocalAndroidPlatforms.findLocalJavaSdk();
                assertEquals(mockSdkDir.getAbsolutePath(), result.getAbsolutePath());
            }
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_AndroidSdkRootValid() throws IOException {
        // Mock environment variable
        File mockSdkDir = tempFolder.newFolder("android-sdk-root");
        File platformsDir = new File(mockSdkDir, "platforms");
        platformsDir.mkdirs();
        
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(mockSdkDir.getAbsolutePath());
            
            // Mock LocalAndroidPlatforms constructor and valid() method
            LocalAndroidPlatforms mockPlatforms = mock(LocalAndroidPlatforms.class);
            when(mockPlatforms.valid()).thenReturn(true);
            
            try (MockedStatic<LocalAndroidPlatforms> mockedLocalAndroidPlatforms = 
                 Mockito.mockStatic(LocalAndroidPlatforms.class)) {
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(any(File.class)))
                    .thenReturn(mockPlatforms);
                
                File result = LocalAndroidPlatforms.findLocalJavaSdk();
                assertEquals(mockSdkDir.getAbsolutePath(), result.getAbsolutePath());
            }
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_AndroidHomeInvalidThenSdkRootValid() throws IOException {
        // Mock environment variables
        File invalidDir = tempFolder.newFolder("invalid-android");
        File validDir = tempFolder.newFolder("valid-android");
        File platformsDir = new File(validDir, "platforms");
        platformsDir.mkdirs();
        
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(invalidDir.getAbsolutePath());
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(validDir.getAbsolutePath());
            
            // Mock LocalAndroidPlatforms constructor and valid() method
            LocalAndroidPlatforms invalidPlatforms = mock(LocalAndroidPlatforms.class);
            LocalAndroidPlatforms validPlatforms = mock(LocalAndroidPlatforms.class);
            when(invalidPlatforms.valid()).thenReturn(false);
            when(validPlatforms.valid()).thenReturn(true);
            
            try (MockedStatic<LocalAndroidPlatforms> mockedLocalAndroidPlatforms = 
                 Mockito.mockStatic(LocalAndroidPlatforms.class)) {
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(eq(invalidDir)))
                    .thenReturn(invalidPlatforms);
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(eq(validDir)))
                    .thenReturn(validPlatforms);
                
                File result = LocalAndroidPlatforms.findLocalJavaSdk();
                assertEquals(validDir.getAbsolutePath(), result.getAbsolutePath());
            }
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_FindInPathWindows() throws IOException {
        // Mock Windows OS
        mockedSystemUtils = Mockito.mockStatic(SystemUtils.class);
        mockedSystemUtils.when(() -> SystemUtils.IS_OS_WINDOWS).thenReturn(true);
        
        // Mock PATH environment variable
        File mockPathDir = tempFolder.newFolder("android-tools");
        File mockPlatformsDir = tempFolder.newFolder("platforms");
        File mockSdkDir = mockPathDir.getParentFile();
        
        // Create mock adb.exe
        File mockAdb = new File(mockPathDir, "adb.exe");
        mockAdb.createNewFile();
        mockAdb.setExecutable(true);
        
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("PATH")).thenReturn(mockPathDir.getAbsolutePath());
            
            // Mock StringUtils.split
            mockedStringUtils = Mockito.mockStatic(StringUtils.class);
            mockedStringUtils.when(() -> StringUtils.split(eq(mockPathDir.getAbsolutePath()), anyChar()))
                .thenReturn(new String[]{mockPathDir.getAbsolutePath()});
            
            // Mock LocalAndroidPlatforms constructor and valid() method
            LocalAndroidPlatforms mockPlatforms = mock(LocalAndroidPlatforms.class);
            when(mockPlatforms.valid()).thenReturn(true);
            
            try (MockedStatic<LocalAndroidPlatforms> mockedLocalAndroidPlatforms = 
                 Mockito.mockStatic(LocalAndroidPlatforms.class)) {
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(eq(mockSdkDir)))
                    .thenReturn(mockPlatforms);
                
                File result = LocalAndroidPlatforms.findLocalJavaSdk();
                assertEquals(mockSdkDir.getAbsolutePath(), result.getAbsolutePath());
            }
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_FindInPathUnix() throws IOException {
        // Mock Unix OS
        mockedSystemUtils = Mockito.mockStatic(SystemUtils.class);
        mockedSystemUtils.when(() -> SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        
        // Mock PATH environment variable
        File mockPathDir = tempFolder.newFolder("android-tools");
        File mockPlatformsDir = tempFolder.newFolder("platforms");
        File mockSdkDir = mockPathDir.getParentFile();
        
        // Create mock adb
        File mockAdb = new File(mockPathDir, "adb");
        mockAdb.createNewFile();
        mockAdb.setExecutable(true);
        
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("PATH")).thenReturn(mockPathDir.getAbsolutePath());
            
            // Mock StringUtils.split
            mockedStringUtils = Mockito.mockStatic(StringUtils.class);
            mockedStringUtils.when(() -> StringUtils.split(eq(mockPathDir.getAbsolutePath()), anyChar()))
                .thenReturn(new String[]{mockPathDir.getAbsolutePath()});
            
            // Mock LocalAndroidPlatforms constructor and valid() method
            LocalAndroidPlatforms mockPlatforms = mock(LocalAndroidPlatforms.class);
            when(mockPlatforms.valid()).thenReturn(true);
            
            try (MockedStatic<LocalAndroidPlatforms> mockedLocalAndroidPlatforms = 
                 Mockito.mockStatic(LocalAndroidPlatforms.class)) {
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(eq(mockSdkDir)))
                    .thenReturn(mockPlatforms);
                
                File result = LocalAndroidPlatforms.findLocalJavaSdk();
                assertEquals(mockSdkDir.getAbsolutePath(), result.getAbsolutePath());
            }
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_NoEnvironmentVariables() throws IOException {
        // Mock no environment variables
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("PATH")).thenReturn("");
            
            // Mock StringUtils.split for empty PATH
            mockedStringUtils = Mockito.mockStatic(StringUtils.class);
            mockedStringUtils.when(() -> StringUtils.split(eq(""), anyChar()))
                .thenReturn(new String[0]);
            
            // Mock Windows OS to test both branches
            mockedSystemUtils = Mockito.mockStatic(SystemUtils.class);
            mockedSystemUtils.when(() -> SystemUtils.IS_OS_WINDOWS).thenReturn(true);
            
            thrown.expect(FileNotFoundException.class);
            thrown.expectMessage("Unable to find the Local Android Java SDK Folder.");
            
            LocalAndroidPlatforms.findLocalJavaSdk();
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_InvalidEnvironmentVariablePath() throws IOException {
        // Mock environment variable with non-existent path
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn("/non/existent/path");
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
            
            // Mock LocalAndroidPlatforms constructor and valid() method
            LocalAndroidPlatforms mockPlatforms = mock(LocalAndroidPlatforms.class);
            when(mockPlatforms.valid()).thenReturn(false);
            
            try (MockedStatic<LocalAndroidPlatforms> mockedLocalAndroidPlatforms = 
                 Mockito.mockStatic(LocalAndroidPlatforms.class)) {
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(any(File.class)))
                    .thenReturn(mockPlatforms);
                
                // Mock Windows OS and empty PATH
                mockedSystemUtils = Mockito.mockStatic(SystemUtils.class);
                mockedSystemUtils.when(() -> SystemUtils.IS_OS_WINDOWS).thenReturn(true);
                mockedSystem.when(() -> System.getenv("PATH")).thenReturn("");
                
                mockedStringUtils = Mockito.mockStatic(StringUtils.class);
                mockedStringUtils.when(() -> StringUtils.split(eq(""), anyChar()))
                    .thenReturn(new String[0]);
                
                thrown.expect(FileNotFoundException.class);
                thrown.expectMessage("Unable to find the Local Android Java SDK Folder.");
                
                LocalAndroidPlatforms.findLocalJavaSdk();
            }
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_ExecutableNotFoundInPath() throws IOException {
        // Mock environment variables not set
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
            
            // Mock PATH with directory but no executables
            File mockPathDir = tempFolder.newFolder("empty-tools");
            mockedSystem.when(() -> System.getenv("PATH")).thenReturn(mockPathDir.getAbsolutePath());
            
            // Mock StringUtils.split
            mockedStringUtils = Mockito.mockStatic(StringUtils.class);
            mockedStringUtils.when(() -> StringUtils.split(eq(mockPathDir.getAbsolutePath()), anyChar()))
                .thenReturn(new String[]{mockPathDir.getAbsolutePath()});
            
            // Mock Windows OS
            mockedSystemUtils = Mockito.mockStatic(SystemUtils.class);
            mockedSystemUtils.when(() -> SystemUtils.IS_OS_WINDOWS).thenReturn(true);
            
            thrown.expect(FileNotFoundException.class);
            thrown.expectMessage("Unable to find the Local Android Java SDK Folder.");
            
            LocalAndroidPlatforms.findLocalJavaSdk();
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_ExecutableFoundButInvalidPlatform() throws IOException {
        // Mock Windows OS
        mockedSystemUtils = Mockito.mockStatic(SystemUtils.class);
        mockedSystemUtils.when(() -> SystemUtils.IS_OS_WINDOWS).thenReturn(true);
        
        // Mock PATH environment variable
        File mockPathDir = tempFolder.newFolder("android-tools");
        File mockSdkDir = mockPathDir.getParentFile();
        
        // Create mock adb.exe
        File mockAdb = new File(mockPathDir, "adb.exe");
        mockAdb.createNewFile();
        mockAdb.setExecutable(true);
        
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("PATH")).thenReturn(mockPathDir.getAbsolutePath());
            
            // Mock StringUtils.split
            mockedStringUtils = Mockito.mockStatic(StringUtils.class);
            mockedStringUtils.when(() -> StringUtils.split(eq(mockPathDir.getAbsolutePath()), anyChar()))
                .thenReturn(new String[]{mockPathDir.getAbsolutePath()});
            
            // Mock LocalAndroidPlatforms constructor and valid() method returning false
            LocalAndroidPlatforms mockPlatforms = mock(LocalAndroidPlatforms.class);
            when(mockPlatforms.valid()).thenReturn(false);
            
            try (MockedStatic<LocalAndroidPlatforms> mockedLocalAndroidPlatforms = 
                 Mockito.mockStatic(LocalAndroidPlatforms.class)) {
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(eq(mockSdkDir)))
                    .thenReturn(mockPlatforms);
                
                thrown.expect(FileNotFoundException.class);
                thrown.expectMessage("Unable to find the Local Android Java SDK Folder.");
                
                LocalAndroidPlatforms.findLocalJavaSdk();
            }
        }
    }
    
    @Test
    public void testFindLocalJavaSdk_MultipleExecutablesInPath() throws IOException {
        // Mock Unix OS
        mockedSystemUtils = Mockito.mockStatic(SystemUtils.class);
        mockedSystemUtils.when(() -> SystemUtils.IS_OS_WINDOWS).thenReturn(false);
        
        // Mock PATH with multiple directories
        File mockPathDir1 = tempFolder.newFolder("tools1");
        File mockPathDir2 = tempFolder.newFolder("tools2");
        File mockSdkDir = mockPathDir2.getParentFile();
        
        // Create mock adb in second directory
        File mockAdb = new File(mockPathDir2, "adb");
        mockAdb.createNewFile();
        mockAdb.setExecutable(true);
        
        try (MockedStatic<System> mockedSystem = Mockito.mockStatic(System.class)) {
            mockedSystem.when(() -> System.getenv("ANDROID_HOME")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("ANDROID_SDK_ROOT")).thenReturn(null);
            mockedSystem.when(() -> System.getenv("PATH"))
                .thenReturn(mockPathDir1.getAbsolutePath() + File.pathSeparator + mockPathDir2.getAbsolutePath());
            
            // Mock StringUtils.split
            mockedStringUtils = Mockito.mockStatic(StringUtils.class);
            mockedStringUtils.when(() -> StringUtils.split(
                eq(mockPathDir1.getAbsolutePath() + File.pathSeparator + mockPathDir2.getAbsolutePath()), 
                anyChar()))
                .thenReturn(new String[]{mockPathDir1.getAbsolutePath(), mockPathDir2.getAbsolutePath()});
            
            // Mock LocalAndroidPlatforms constructor and valid() method
            LocalAndroidPlatforms mockPlatforms = mock(LocalAndroidPlatforms.class);
            when(mockPlatforms.valid()).thenReturn(true);
            
            try (MockedStatic<LocalAndroidPlatforms> mockedLocalAndroidPlatforms = 
                 Mockito.mockStatic(LocalAndroidPlatforms.class)) {
                mockedLocalAndroidPlatforms.when(() -> new LocalAndroidPlatforms(eq(mockSdkDir)))
                    .thenReturn(mockPlatforms);
                
                File result = LocalAndroidPlatforms.findLocalJavaSdk();
                assertEquals(mockSdkDir.getAbsolutePath(), result.getAbsolutePath());
            }
        }
    }
}