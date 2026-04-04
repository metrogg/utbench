import static org.mockito.Matchers.any;
import static org.mockito.Matchers.anyInt;
import static org.mockito.Matchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;
import static org.powermock.api.mockito.PowerMockito.doNothing;
import static org.powermock.api.mockito.PowerMockito.mockStatic;
import static org.powermock.api.mockito.PowerMockito.whenNew;

import java.io.IOException;
import java.io.InputStream;

import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import jxl.Cell;
import jxl.Sheet;
import jxl.Workbook;
import jxl.read.biff.BiffException;

@RunWith(PowerMockRunner.class)
@PrepareForTest({ Workbook.class, StringUtils.class, Debug.class })
public class XlsLegendParserTest {

    private XlsLegendParser parser;
    private InputStream mockInputStream;
    private Workbook mockWorkbook;
    private Sheet mockSheet;

    @Before
    public void setUp() throws Exception {
        parser = new XlsLegendParser();
        mockInputStream = mock(InputStream.class);
        mockWorkbook = mock(Workbook.class);
        mockSheet = mock(Sheet.class);

        // Mock static dependencies
        mockStatic(Workbook.class);
        mockStatic(StringUtils.class);
        mockStatic(Debug.class);

        // Mock Debug.trace to do nothing
        doNothing().when(Debug.class);
        Debug.trace(any(Throwable.class));

        // Default behavior for StringUtils.isNullOrEmpty to be false (content exists)
        // Specific tests will override this where necessary
        when(StringUtils.isNullOrEmpty(anyString())).thenReturn(false);

        // Default Workbook behavior
        when(Workbook.getWorkbook(any(InputStream.class))).thenReturn(mockWorkbook);
    }

    @Test
    public void testParseGlobalSuccess() throws Exception {
        // Setup Sheet
        when(mockWorkbook.getSheet(anyString())).thenReturn(mockSheet);
        when(mockSheet.getRows()).thenReturn(2); // Header + 1 Data Row

        // Mock Headers
        Cell valHeader = mockCell(0, "Value");
        Cell lblHeader = mockCell(1, "Label");
        Cell redHeader = mockCell(2, "Red");
        Cell grnHeader = mockCell(3, "Green");
        Cell bluHeader = mockCell(4, "Blue");

        when(mockSheet.findCell("Value")).thenReturn(valHeader);
        when(mockSheet.findCell("Label")).thenReturn(lblHeader);
        when(mockSheet.findCell("Red")).thenReturn(redHeader);
        when(mockSheet.findCell("Green")).thenReturn(grnHeader);
        when(mockSheet.findCell("Blue")).thenReturn(bluHeader);

        // Mock Data Row (Index 1)
        when(mockSheet.getCell(0, 1)).thenReturn(mockCell(0, "100"));
        when(mockSheet.getCell(1, 1)).thenReturn(mockCell(1, "Description"));
        when(mockSheet.getCell(2, 1)).thenReturn(mockCell(2, "255"));
        when(mockSheet.getCell(3, 1)).thenReturn(mockCell(3, "0"));
        when(mockSheet.getCell(4, 1)).thenReturn(mockCell(4, "0"));

        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);

        // Assert
        assertEquals(1, result.length);
        assertEquals(100, result[0].getValue());
        assertEquals("Class_1", result[0].getName());
        assertEquals("Description", result[0].getDescr());
        assertEquals(255, result[0].getColor().getRed());
        assertEquals(0, result[0].getColor().getGreen());
        assertEquals(0, result[0].getColor().getBlue());
    }

    @Test(expected = IllegalArgumentException.class)
    public void testParseRegionalSheetMissing() throws Exception {
        // Setup: Regional sheet requested, but returns null
        when(mockWorkbook.getSheet(anyString())).thenReturn(null);

        // Execute
        parser.parse(mockInputStream, true);
    }

    @Test
    public void testParseSkipEmptyValueRow() throws Exception {
        when(mockWorkbook.getSheet(anyString())).thenReturn(mockSheet);
        when(mockSheet.getRows()).thenReturn(2);

        // Mock Headers
        setupHeaders(mockSheet);

        // Mock Data Row with (effectively) empty value
        when(StringUtils.isNullOrEmpty("")).thenReturn(true);
        
        when(mockSheet.getCell(0, 1)).thenReturn(mockCell(0, ""));
        when(mockSheet.getCell(1, 1)).thenReturn(mockCell(1, "Desc"));
        when(mockSheet.getCell(2, 1)).thenReturn(mockCell(2, "10"));
        when(mockSheet.getCell(3, 1)).thenReturn(mockCell(3, "10"));
        when(mockSheet.getCell(4, 1)).thenReturn(mockCell(4, "10"));

        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);

        // Assert: Array is created but slot is not filled
        assertEquals(1, result.length);
        assertNull(result[0]);
    }

    @Test
    public void testParseHandlesIOException() throws Exception {
        // Setup: Workbook.getWorkbook throws IOException
        when(Workbook.getWorkbook(any(InputStream.class))).thenThrow(new IOException("Test IO Error"));

        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);

        // Assert: Should return empty array and not crash
        assertEquals(0, result.length);
    }

    @Test
    public void testParseHandlesBiffException() throws Exception {
        // Setup: Workbook.getWorkbook throws BiffException
        when(Workbook.getWorkbook(any(InputStream.class))).thenThrow(new BiffException("Test Biff Error"));

        // Execute
        LegendClass[] result = parser.parse(mockInputStream, false);

        // Assert: Should return empty array and not crash
        assertEquals(0, result.length);
    }
    
    @Test
    public void testParseRegionalSuccess() throws Exception {
        // Setup Sheet for Regional
        when(mockWorkbook.getSheet(anyString())).thenReturn(mockSheet);
        when(mockSheet.getRows()).thenReturn(2);

        setupHeaders(mockSheet);

        // Mock Data
        when(mockSheet.getCell(0, 1)).thenReturn(mockCell(0, "5"));
        when(mockSheet.getCell(1, 1)).thenReturn(mockCell(1, "Regional Class"));
        when(mockSheet.getCell(2, 1)).thenReturn(mockCell(2, "0"));
        when(mockSheet.getCell(3, 1)).thenReturn(mockCell(3, "255"));
        when(mockSheet.getCell(4, 1)).thenReturn(mockCell(4, "0"));

        // Execute
        LegendClass[] result = parser.parse(mockInputStream, true);

        // Assert
        assertEquals(1, result.length);
        assertEquals(5, result[0].getValue());
        assertEquals("Regional Class", result[0].getDescr());
    }

    // Helper method to reduce boilerplate
    private void setupHeaders(Sheet sheet) {
        when(sheet.findCell("Value")).thenReturn(mockCell(0, "Value"));
        when(sheet.findCell("Label")).thenReturn(mockCell(1, "Label"));
        when(sheet.findCell("Red")).thenReturn(mockCell(2, "Red"));
        when(sheet.findCell("Green")).thenReturn(mockCell(3, "Green"));
        when(sheet.findCell("Blue")).thenReturn(mockCell(4, "Blue"));
    }

    private Cell mockCell(int col, String content) {
        Cell cell = mock(Cell.class);
        when(cell.getColumn()).thenReturn(col);
        when(cell.getContents()).thenReturn(content);
        return cell;
    }
}