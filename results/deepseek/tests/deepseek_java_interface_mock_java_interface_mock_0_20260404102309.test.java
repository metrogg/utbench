```java
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.MockitoJUnitRunner;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import jxl.Cell;
import jxl.Sheet;
import jxl.Workbook;
import jxl.read.biff.BiffException;
import static org.junit.Assert.*;
import static org.mockito.Mockito.*;

@RunWith(MockitoJUnitRunner.class)
public class XlsLegendParserTest {

    private XlsLegendParser parser;

    @Mock
    private Workbook mockWorkbook;

    @Mock
    private Sheet mockSheet;

    @Mock
    private Cell mockValueCell;

    @Mock
    private Cell mockLabelCell;

    @Mock
    private Cell mockRedCell;

    @Mock
    private Cell mockGreenCell;

    @Mock
    private Cell mockBlueCell;

    @Before
    public void setUp() {
        parser = new XlsLegendParser();
    }

    @Test
    public void testParseGlobalLegendWithValidData() throws Exception {
        // Setup mock workbook and sheet
        when(mockWorkbook.getSheet("Global")).thenReturn(mockSheet);
        when(mockSheet.findCell("Value")).thenReturn(mockValueCell);
        when(mockSheet.findCell("Label")).thenReturn(mockLabelCell);
        when(mockSheet.findCell("Red")).thenReturn(mockRedCell);
        when(mockSheet.findCell("Green")).thenReturn(mockGreenCell);
        when(mockSheet.findCell("Blue")).thenReturn(mockBlueCell);
        
        when(mockValueCell.getColumn()).thenReturn(0);
        when(mockLabelCell.getColumn()).thenReturn(1);
        when(mockRedCell.getColumn()).thenReturn(2);
        when(mockGreenCell.getColumn()).thenReturn(3);
        when(mockBlueCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(3);
        
        // Row 1 data
        Cell mockValueCell1 = mock(Cell.class);
        Cell mockLabelCell1 = mock(Cell.class);
        Cell mockRedCell1 = mock(Cell.class);
        Cell mockGreenCell1 = mock(Cell.class);
        Cell mockBlueCell1 = mock(Cell.class);
        
        when(mockSheet.getCell(0, 1)).thenReturn(mockValueCell1);
        when(mockSheet.getCell(1, 1)).thenReturn(mockLabelCell1);
        when(mockSheet.getCell(2, 1)).thenReturn(mockRedCell1);
        when(mockSheet.getCell(3, 1)).thenReturn(mockGreenCell1);
        when(mockSheet.getCell(4, 1)).thenReturn(mockBlueCell1);
        
        when(mockValueCell1.getContents()).thenReturn("1");
        when(mockLabelCell1.getContents()).thenReturn("Forest");
        when(mockRedCell1.getContents()).thenReturn("0");
        when(mockGreenCell1.getContents()).thenReturn("255");
        when(mockBlueCell1.getContents()).thenReturn("0");
        
        // Row 2 data
        Cell mockValueCell2 = mock(Cell.class);
        Cell mockLabelCell2 = mock(Cell.class);
        Cell mockRedCell2 = mock(Cell.class);
        Cell mockGreenCell2 = mock(Cell.class);
        Cell mockBlueCell2 = mock(Cell.class);
        
        when(mockSheet.getCell(0, 2)).thenReturn(mockValueCell2);
        when(mockSheet.getCell(1, 2)).thenReturn(mockLabelCell2);
        when(mockSheet.getCell(2, 2)).thenReturn(mockRedCell2);
        when(mockSheet.getCell(3, 2)).thenReturn(mockGreenCell2);
        when(mockSheet.getCell(4, 2)).thenReturn(mockBlueCell2);
        
        when(mockValueCell2.getContents()).thenReturn("2");
        when(mockLabelCell2.getContents()).thenReturn("Water");
        when(mockRedCell2.getContents()).thenReturn("0");
        when(mockGreenCell2.getContents()).thenReturn("0");
        when(mockBlueCell2.getContents()).thenReturn("255");
        
        // Create test input stream
        InputStream mockInputStream = mock(InputStream.class);
        
        // Mock Workbook factory method
        Workbook originalWorkbook = Workbook.getWorkbook(mockInputStream);
        try {
            java.lang.reflect.Field field = Workbook.class.getDeclaredField("testWorkbook");
            field.setAccessible(true);
            field.set(null, mockWorkbook);
        } catch (Exception e) {
            // Reflection fallback - use dependency injection pattern
        }
        
        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        // Verify
        assertNotNull(result);
        assertEquals(2, result.length);
        
        assertEquals(1, result[0].getValue());
        assertEquals("Forest", result[0].getDescription());
        assertEquals("Class_1", result[0].getName());
        assertEquals(0, result[0].getColor().getRed());
        assertEquals(255, result[0].getColor().getGreen());
        assertEquals(0, result[0].getColor().getBlue());
        
        assertEquals(2, result[1].getValue());
        assertEquals("Water", result[1].getDescription());
        assertEquals("Class_2", result[1].getName());
        assertEquals(0, result[1].getColor().getRed());
        assertEquals(0, result[1].getColor().getGreen());
        assertEquals(255, result[1].getColor().getBlue());
        
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseRegionalLegendWithValidData() throws Exception {
        // Setup mock workbook and sheet
        when(mockWorkbook.getSheet("Regional")).thenReturn(mockSheet);
        when(mockSheet.findCell("Value")).thenReturn(mockValueCell);
        when(mockSheet.findCell("Label")).thenReturn(mockLabelCell);
        when(mockSheet.findCell("Red")).thenReturn(mockRedCell);
        when(mockSheet.findCell("Green")).thenReturn(mockGreenCell);
        when(mockSheet.findCell("Blue")).thenReturn(mockBlueCell);
        
        when(mockValueCell.getColumn()).thenReturn(0);
        when(mockLabelCell.getColumn()).thenReturn(1);
        when(mockRedCell.getColumn()).thenReturn(2);
        when(mockGreenCell.getColumn()).thenReturn(3);
        when(mockBlueCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(2);
        
        // Row 1 data
        Cell mockValueCell1 = mock(Cell.class);
        Cell mockLabelCell1 = mock(Cell.class);
        Cell mockRedCell1 = mock(Cell.class);
        Cell mockGreenCell1 = mock(Cell.class);
        Cell mockBlueCell1 = mock(Cell.class);
        
        when(mockSheet.getCell(0, 1)).thenReturn(mockValueCell1);
        when(mockSheet.getCell(1, 1)).thenReturn(mockLabelCell1);
        when(mockSheet.getCell(2, 1)).thenReturn(mockRedCell1);
        when(mockSheet.getCell(3, 1)).thenReturn(mockGreenCell1);
        when(mockSheet.getCell(4, 1)).thenReturn(mockBlueCell1);
        
        when(mockValueCell1.getContents()).thenReturn("10");
        when(mockLabelCell1.getContents()).thenReturn("Urban");
        when(mockRedCell1.getContents()).thenReturn("255");
        when(mockGreenCell1.getContents()).thenReturn("0");
        when(mockBlueCell1.getContents()).thenReturn("0");
        
        // Create test input stream
        InputStream mockInputStream = mock(InputStream.class);
        
        // Mock Workbook factory method
        Workbook originalWorkbook = Workbook.getWorkbook(mockInputStream);
        try {
            java.lang.reflect.Field field = Workbook.class.getDeclaredField("testWorkbook");
            field.setAccessible(true);
            field.set(null, mockWorkbook);
        } catch (Exception e) {
            // Reflection fallback
        }
        
        // Execute
        LegendClass[] result = parser.parse(mockInputStream, true);
        
        // Verify
        assertNotNull(result);
        assertEquals(1, result.length);
        
        assertEquals(10, result[0].getValue());
        assertEquals("Urban", result[0].getDescription());
        assertEquals("Class_1", result[0].getName());
        assertEquals(255, result[0].getColor().getRed());
        assertEquals(0, result[0].getColor().getGreen());
        assertEquals(0, result[0].getColor().getBlue());
        
        verify(mockWorkbook).close();
    }

    @Test(expected = IllegalArgumentException.class)
    public void testParseRegionalLegendWithMissingSheet() throws Exception {
        // Setup mock workbook without regional sheet
        when(mockWorkbook.getSheet("Regional")).thenReturn(null);
        
        // Create test input stream
        InputStream mockInputStream = mock(InputStream.class);
        
        // Mock Workbook factory method
        Workbook originalWorkbook = Workbook.getWorkbook(mockInputStream);
        try {
            java.lang.reflect.Field field = Workbook.class.getDeclaredField("testWorkbook");
            field.setAccessible(true);
            field.set(null, mockWorkbook);
        } catch (Exception e) {
            // Reflection fallback
        }
        
        // Execute - should throw IllegalArgumentException
        parser.parse(mockInputStream, true);
    }

    @Test
    public void testParseWithEmptyValueCell() throws Exception {
        // Setup mock workbook and sheet
        when(mockWorkbook.getSheet("Global")).thenReturn(mockSheet);
        when(mockSheet.findCell("Value")).thenReturn(mockValueCell);
        when(mockSheet.findCell("Label")).thenReturn(mockLabelCell);
        when(mockSheet.findCell("Red")).thenReturn(mockRedCell);
        when(mockSheet.findCell("Green")).thenReturn(mockGreenCell);
        when(mockSheet.findCell("Blue")).thenReturn(mockBlueCell);
        
        when(mockValueCell.getColumn()).thenReturn(0);
        when(mockLabelCell.getColumn()).thenReturn(1);
        when(mockRedCell.getColumn()).thenReturn(2);
        when(mockGreenCell.getColumn()).thenReturn(3);
        when(mockBlueCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(3);
        
        // Row 1 - empty value cell (should be skipped)
        Cell mockValueCell1 = mock(Cell.class);
        when(mockSheet.getCell(0, 1)).thenReturn(mockValueCell1);
        when(mockValueCell1.getContents()).thenReturn("");
        
        // Row 2 - valid data
        Cell mockValueCell2 = mock(Cell.class);
        Cell mockLabelCell2 = mock(Cell.class);
        Cell mockRedCell2 = mock(Cell.class);
        Cell mockGreenCell2 = mock(Cell.class);
        Cell mockBlueCell2 = mock(Cell.class);
        
        when(mockSheet.getCell(0, 2)).thenReturn(mockValueCell2);
        when(mockSheet.getCell(1, 2)).thenReturn(mockLabelCell2);
        when(mockSheet.getCell(2, 2)).thenReturn(mockRedCell2);
        when(mockSheet.getCell(3, 2)).thenReturn(mockGreenCell2);
        when(mockSheet.getCell(4, 2)).thenReturn(mockBlueCell2);
        
        when(mockValueCell2.getContents()).thenReturn("5");
        when(mockLabelCell2.getContents()).thenReturn("Desert");
        when(mockRedCell2.getContents()).thenReturn("255");
        when(mockGreenCell2.getContents()).thenReturn("255");
        when(mockBlueCell2.getContents()).thenReturn("0");
        
        // Create test input stream
        InputStream mockInputStream = mock(InputStream.class);
        
        // Mock Workbook factory method
        Workbook originalWorkbook = Workbook.getWorkbook(mockInputStream);
        try {
            java.lang.reflect.Field field = Workbook.class.getDeclaredField("testWorkbook");
            field.setAccessible(true);
            field.set(null, mockWorkbook);
        } catch (Exception e) {
            // Reflection fallback
        }
        
        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        // Verify - only one valid row
        assertNotNull(result);
        assertEquals(2, result.length); // Array size is rows-1 = 2
        assertNull(result[0]); // First element should be null (skipped row)
        assertNotNull(result[1]); // Second element should have data
        
        assertEquals(5, result[1].getValue());
        assertEquals("Desert", result[1].getDescription());
        assertEquals("Class_2", result[1].getName());
        
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseWithSingleRowSheet() throws Exception {
        // Setup mock workbook and sheet with only header row
        when(mockWorkbook.getSheet("Global")).thenReturn(mockSheet);
        when(mockSheet.findCell("Value")).thenReturn(mockValueCell);
        when(mockSheet.findCell("Label")).thenReturn(mockLabelCell);
        when(mockSheet.findCell("Red")).thenReturn(mockRedCell);
        when(mockSheet.findCell("Green")).thenReturn(mockGreenCell);
        when(mockSheet.findCell("Blue")).thenReturn(mockBlueCell);
        
        when(mockValueCell.getColumn()).thenReturn(0);
        when(mockLabelCell.getColumn()).thenReturn(1);
        when(mockRedCell.getColumn()).thenReturn(2);
        when(mockGreenCell.getColumn()).thenReturn(3);
        when(mockBlueCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(1); // Only header row
        
        // Create test input stream
        InputStream mockInputStream = mock(InputStream.class);
        
        // Mock Workbook factory method
        Workbook originalWorkbook = Workbook.getWorkbook(mockInputStream);
        try {
            java.lang.reflect.Field field = Workbook.class.getDeclaredField("testWorkbook");
            field.setAccessible(true);
            field.set(null, mockWorkbook);
        } catch (Exception e) {
            // Reflection fallback
        }
        
        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        // Verify - empty array
        assertNotNull(result);
        assertEquals(0, result.length);
        
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseWithNullInputStream() {
        // Execute with null input stream
        LegendClass[] result = parser.parse(null, false);
        
        // Verify - empty array returned (due to exception handling)
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    public void testParseWithBiffException() throws Exception {
        // Setup to throw BiffException when getting workbook
        InputStream mockInputStream = mock(InputStream.class);
        when(Workbook.getWorkbook(mockInputStream)).thenThrow(new BiffException("Invalid Excel format"));
        
        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        // Verify - empty array returned
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    public void testParseWithIOException() throws Exception {
        // Setup to throw IOException when getting workbook
        InputStream mockInputStream = mock(InputStream.class);
        when(Workbook.getWorkbook(mockInputStream)).thenThrow(new IOException("Stream error"));
        
        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        // Verify - empty array returned
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    public void testParseWithInvalidNumberFormat() throws Exception {
        // Setup mock workbook and sheet
        when(mockWorkbook.getSheet("Global")).thenReturn(mockSheet);
        when(mockSheet.findCell("Value")).thenReturn(mockValueCell);
        when(mockSheet.findCell("Label")).thenReturn(mockLabelCell);
        when(mockSheet.findCell("Red")).thenReturn(mockRedCell);
        when(mockSheet.findCell("Green")).thenReturn(mockGreenCell);
        when(mockSheet.findCell("Blue")).thenReturn(mockBlueCell);
        
        when(mockValueCell.getColumn()).thenReturn(0);
        when(mockLabelCell.getColumn()).thenReturn(1);
        when(mockRedCell.getColumn()).thenReturn(2);
        when(mockGreenCell.getColumn()).thenReturn(3);
        when(mockBlueCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(2);
        
        // Row with invalid number format
        Cell mockValueCell1 = mock(Cell.class);
        when(mockSheet.getCell(0, 1)).thenReturn(mockValueCell1);
        when(mockValueCell1.getContents()).thenReturn("not-a-number");
        
        // Create test input stream
        InputStream mockInputStream = mock(InputStream.class);
        
        // Mock Workbook factory method
        Workbook originalWorkbook = Workbook.getWorkbook(mockInputStream);
        try {
            java.lang.reflect.Field field = Workbook.class.getDeclaredField("testWorkbook");
            field.setAccessible(true);
            field.set(null, mockWorkbook);
        } catch (Exception e) {
            // Reflection fallback
        }
        
        // Execute - should handle NumberFormatException internally
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        // Verify - empty array due to exception
        assertNotNull(result);
        assertEquals(1, result.length); // Array size is rows-1 = 1
        assertNull(result[0]); // Element should be null due to exception
        
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseWithTrimmedLabel() throws Exception {
        // Setup mock workbook and sheet
        when(mockWorkbook.getSheet("Global")).thenReturn(mockSheet);
        when(mockSheet.findCell("Value")).thenReturn(mockValueCell);
        when(mockSheet.findCell("Label")).thenReturn(mockLabelCell);
        when(mockSheet.findCell("Red")).thenReturn(mockRedCell);
        when(mockSheet.findCell("Green")).thenReturn(mockGreenCell);
        when(mockSheet.findCell("Blue")).thenReturn(mockBlueCell);
        
        when(mockValueCell.getColumn()).thenReturn(0);
        when(m