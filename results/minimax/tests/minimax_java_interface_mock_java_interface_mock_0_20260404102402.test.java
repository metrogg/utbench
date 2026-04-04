import org.junit.Before;
import org.junit.Test;
import org.mockito.ArgumentCaptor;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import static org.mockito.Mockito.*;

import java.io.ByteArrayInputStream;
import java.io.InputStream;

import jxl.Cell;
import jxl.Sheet;
import jxl.Workbook;
import jxl.read.biff.BiffException;

public class XlsLegendParserTest {

    private XlsLegendParser parser;

    @Mock
    private InputStream mockInputStream;

    @Mock
    private Workbook mockWorkbook;

    @Mock
    private Sheet mockSheet;

    @Mock
    private Cell mockCell;

    @Before
    public void setUp() {
        MockitoAnnotations.initMocks(this);
        parser = new XlsLegendParser();
    }

    @Test
    public void testGlobalLegendParsingRegionalSheet() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
        when(mockWorkbook.getSheet("Legend")).thenReturn(null);

        try {
            parser.parse(mockInputStream, true);
        } catch (IllegalArgumentException e) {
            assertEquals("Given xls-legend is not regional.", e.getMessage());
        }
    }

    @Test
    public void testParseGlobalLegendSuccess() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
        when(mockWorkbook.getSheet("Global_Legend")).thenReturn(mockSheet);
        
        when(mockSheet.findCell("VALUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(0);
        when(mockSheet.findCell("LABEL")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(1);
        when(mockSheet.findCell("RED")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(2);
        when(mockSheet.findCell("GREEN")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(3);
        when(mockSheet.findCell("BLUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(3);
        
        Cell mockValueCell = mock(Cell.class);
        Cell mockLabelCell = mock(Cell.class);
        Cell mockRedCell = mock(Cell.class);
        Cell mockGreenCell = mock(Cell.class);
        Cell mockBlueCell = mock(Cell.class);
        
        when(mockSheet.getCell(0, 1)).thenReturn(mockValueCell);
        when(mockSheet.getCell(1, 1)).thenReturn(mockLabelCell);
        when(mockSheet.getCell(2, 1)).thenReturn(mockRedCell);
        when(mockSheet.getCell(3, 1)).thenReturn(mockGreenCell);
        when(mockSheet.getCell(4, 1)).thenReturn(mockBlueCell);
        
        when(mockValueCell.getContents()).thenReturn("1");
        when(mockLabelCell.getContents()).thenReturn("  Class One  ");
        when(mockRedCell.getContents()).thenReturn("255");
        when(mockGreenCell.getContents()).thenReturn("128");
        when(mockBlueCell.getContents()).thenReturn("0");
        
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
        when(mockLabelCell2.getContents()).thenReturn("Class Two");
        when(mockRedCell2.getContents()).thenReturn("0");
        when(mockGreenCell2.getContents()).thenReturn("255");
        when(mockBlueCell2.getContents()).thenReturn("0");
        
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        assertEquals(2, result.length);
        assertEquals(1, result[0].getValue());
        assertEquals("Class One", result[0].getDescription().trim());
        assertEquals("Class_1", result[0].getName());
        assertEquals(255, result[0].getColor().getRed());
        assertEquals(128, result[0].getColor().getGreen());
        assertEquals(0, result[0].getColor().getBlue());
        
        assertEquals(2, result[1].getValue());
        assertEquals("Class Two", result[1].getDescription());
        
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseRegionalLegendSuccess() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
        when(mockWorkbook.getSheet("Legend")).thenReturn(mockSheet);
        
        when(mockSheet.findCell("VALUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(0);
        when(mockSheet.findCell("LABEL")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(1);
        when(mockSheet.findCell("RED")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(2);
        when(mockSheet.findCell("GREEN")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(3);
        when(mockSheet.findCell("BLUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(2);
        
        Cell mockValueCell = mock(Cell.class);
        Cell mockLabelCell = mock(Cell.class);
        Cell mockRedCell = mock(Cell.class);
        Cell mockGreenCell = mock(Cell.class);
        Cell mockBlueCell = mock(Cell.class);
        
        when(mockSheet.getCell(0, 1)).thenReturn(mockValueCell);
        when(mockSheet.getCell(1, 1)).thenReturn(mockLabelCell);
        when(mockSheet.getCell(2, 1)).thenReturn(mockRedCell);
        when(mockSheet.getCell(3, 1)).thenReturn(mockGreenCell);
        when(mockSheet.getCell(4, 1)).thenReturn(mockBlueCell);
        
        when(mockValueCell.getContents()).thenReturn("10");
        when(mockLabelCell.getContents()).thenReturn("Regional");
        when(mockRedCell.getContents()).thenReturn("100");
        when(mockGreenCell.getContents()).thenReturn("100");
        when(mockBlueCell.getContents()).thenReturn("100");
        
        LegendClass[] result = parser.parse(mockInputStream, true);
        
        assertEquals(1, result.length);
        assertEquals(10, result[0].getValue());
        assertEquals("Regional", result[0].getDescription());
        
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseEmptyValueCell() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
        when(mockWorkbook.getSheet("Global_Legend")).thenReturn(mockSheet);
        
        when(mockSheet.findCell("VALUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(0);
        when(mockSheet.findCell("LABEL")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(1);
        when(mockSheet.findCell("RED")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(2);
        when(mockSheet.findCell("GREEN")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(3);
        when(mockSheet.findCell("BLUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(3);
        
        Cell mockValueCell = mock(Cell.class);
        when(mockSheet.getCell(0, 1)).thenReturn(mockValueCell);
        when(mockValueCell.getContents()).thenReturn("");
        
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        assertEquals(0, result.length);
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseNullValueCell() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
        when(mockWorkbook.getSheet("Global_Legend")).thenReturn(mockSheet);
        
        when(mockSheet.findCell("VALUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(0);
        when(mockSheet.findCell("LABEL")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(1);
        when(mockSheet.findCell("RED")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(2);
        when(mockSheet.findCell("GREEN")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(3);
        when(mockSheet.findCell("BLUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(2);
        
        Cell mockValueCell = mock(Cell.class);
        when(mockSheet.getCell(0, 1)).thenReturn(mockValueCell);
        when(mockValueCell.getContents()).thenReturn(null);
        
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        assertEquals(0, result.length);
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseSheetWithOnlyHeader() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
        when(mockWorkbook.getSheet("Global_Legend")).thenReturn(mockSheet);
        
        when(mockSheet.findCell("VALUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(0);
        when(mockSheet.findCell("LABEL")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(1);
        when(mockSheet.findCell("RED")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(2);
        when(mockSheet.findCell("GREEN")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(3);
        when(mockSheet.findCell("BLUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(1);
        
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        assertEquals(0, result.length);
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseWorkbookClosesOnBiffException() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenThrow(new BiffException("test"));
        
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        assertEquals(0, result.length);
    }

    @Test
    public void testParseWorkbookClosesOnIOException() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenThrow(new java.io.IOException("test"));
        
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        assertEquals(0, result.length);
    }

    @Test
    public void testParseWorkbookIsClosedEvenOnException() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
        when(mockWorkbook.getSheet("Global_Legend")).thenReturn(mockSheet);
        
        when(mockSheet.findCell("VALUE")).thenThrow(new RuntimeException("Test exception"));
        
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        assertEquals(0, result.length);
        verify(mockWorkbook).close();
    }

    @Test
    public void testParseWithMixedEmptyAndValidRows() throws Exception {
        when(Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
        when(mockWorkbook.getSheet("Global_Legend")).thenReturn(mockSheet);
        
        when(mockSheet.findCell("VALUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(0);
        when(mockSheet.findCell("LABEL")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(1);
        when(mockSheet.findCell("RED")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(2);
        when(mockSheet.findCell("GREEN")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(3);
        when(mockSheet.findCell("BLUE")).thenReturn(mockCell);
        when(mockCell.getColumn()).thenReturn(4);
        
        when(mockSheet.getRows()).thenReturn(4);
        
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
        when(mockLabelCell1.getContents()).thenReturn("Valid");
        when(mockRedCell1.getContents()).thenReturn("0");
        when(mockGreenCell1.getContents()).thenReturn("0");
        when(mockBlueCell1.getContents()).thenReturn("0");
        
        Cell mockValueCell2 = mock(Cell.class);
        when(mockSheet.getCell(0, 2)).thenReturn(mockValueCell2);
        when(mockValueCell2.getContents()).thenReturn("");
        
        Cell mockValueCell3 = mock(Cell.class);
        Cell mockLabelCell3 = mock(Cell.class);
        Cell mockRedCell3 = mock(Cell.class);
        Cell mockGreenCell3 = mock(Cell.class);
        Cell mockBlueCell3 = mock(Cell.class);
        
        when(mockSheet.getCell(0, 3)).thenReturn(mockValueCell3);
        when(mockSheet.getCell(1, 3)).thenReturn(mockLabelCell3);
        when(mockSheet.getCell(2, 3)).thenReturn(mockRedCell3);
        when(mockSheet.getCell(3, 3)).thenReturn(mockGreenCell3);
        when(mockSheet.getCell(4, 3)).thenReturn(mockBlueCell3);
        
        when(mockValueCell3.getContents()).thenReturn("2");
        when(mockLabelCell3.getContents()).thenReturn("Also Valid");
        when(mockRedCell3.getContents()).thenReturn("255");
        when(mockGreenCell3.getContents()).thenReturn("255");
        when(mockBlueCell3.getContents()).thenReturn("255");
        
        LegendClass[] result = parser.parse(mockInputStream, false);
        
        assertEquals(2, result.length);
        assertEquals(1, result[0].getValue());
        assertEquals(2, result[1].getValue());
        
        verify(mockWorkbook).close();
    }
}