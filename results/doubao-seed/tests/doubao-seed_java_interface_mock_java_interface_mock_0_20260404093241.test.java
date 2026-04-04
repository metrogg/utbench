import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.MockedStatic;
import org.mockito.junit.MockitoJUnitRunner;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.awt.Color;
import jxl.Workbook;
import jxl.Sheet;
import jxl.Cell;
import jxl.read.biff.BiffException;
import static org.junit.Assert.*;
import static org.mockito.Mockito.*;

@RunWith(MockitoJUnitRunner.class)
public class XlsLegendParserTest {

    private XlsLegendParser parser;

    @Before
    public void setUp() {
        parser = new XlsLegendParser();
    }

    @Test(expected = IllegalArgumentException.class)
    public void testGlobalLegendParsingRegionalSheet() throws Exception {
        InputStream mockIs = new ByteArrayInputStream(new byte[0]);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            Workbook mockWorkbook = mock(Workbook.class);
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenReturn(mockWorkbook);
            when(mockWorkbook.getSheet(XlsLegendParser.REGIONAL_LEGEND_SHEET)).thenReturn(null);
            parser.parse(mockIs, true);
        }
    }

    @Test
    public void testParseRegionalLegendSuccess() throws Exception {
        InputStream mockIs = new ByteArrayInputStream(new byte[0]);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            Workbook mockWorkbook = mock(Workbook.class);
            Sheet mockSheet = mock(Sheet.class);
            Cell valueHeader = mock(Cell.class);
            Cell labelHeader = mock(Cell.class);
            Cell redHeader = mock(Cell.class);
            Cell greenHeader = mock(Cell.class);
            Cell blueHeader = mock(Cell.class);
            Cell valueCell = mock(Cell.class);
            Cell labelCell = mock(Cell.class);
            Cell redCell = mock(Cell.class);
            Cell greenCell = mock(Cell.class);
            Cell blueCell = mock(Cell.class);

            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenReturn(mockWorkbook);
            when(mockWorkbook.getSheet(XlsLegendParser.REGIONAL_LEGEND_SHEET)).thenReturn(mockSheet);

            when(mockSheet.findCell(XlsLegendParser.VALUE)).thenReturn(valueHeader);
            when(valueHeader.getColumn()).thenReturn(0);
            when(mockSheet.findCell(XlsLegendParser.LABEL)).thenReturn(labelHeader);
            when(labelHeader.getColumn()).thenReturn(1);
            when(mockSheet.findCell(XlsLegendParser.RED)).thenReturn(redHeader);
            when(redHeader.getColumn()).thenReturn(2);
            when(mockSheet.findCell(XlsLegendParser.GREEN)).thenReturn(greenHeader);
            when(greenHeader.getColumn()).thenReturn(3);
            when(mockSheet.findCell(XlsLegendParser.BLUE)).thenReturn(blueHeader);
            when(blueHeader.getColumn()).thenReturn(4);

            when(mockSheet.getRows()).thenReturn(2);
            when(mockSheet.getCell(0, 1)).thenReturn(valueCell);
            when(valueCell.getContents()).thenReturn("101");
            when(mockSheet.getCell(1, 1)).thenReturn(labelCell);
            when(labelCell.getContents()).thenReturn("  Forest Area  ");
            when(mockSheet.getCell(2, 1)).thenReturn(redCell);
            when(redCell.getContents()).thenReturn("0");
            when(mockSheet.getCell(3, 1)).thenReturn(greenCell);
            when(greenCell.getContents()).thenReturn("128");
            when(mockSheet.getCell(4, 1)).thenReturn(blueCell);
            when(blueCell.getContents()).thenReturn("0");

            LegendClass[] result = parser.parse(mockIs, true);

            assertEquals(1, result.length);
            assertEquals(101, result[0].getValue());
            assertEquals("Forest Area", result[0].getDescr());
            assertEquals("Class_1", result[0].getName());
            assertEquals(new Color(0, 128, 0), result[0].getColor());
            verify(mockWorkbook, times(1)).close();
        }
    }

    @Test
    public void testParseGlobalLegendSuccess() throws Exception {
        InputStream mockIs = new ByteArrayInputStream(new byte[0]);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            Workbook mockWorkbook = mock(Workbook.class);
            Sheet mockSheet = mock(Sheet.class);
            Cell valueHeader = mock(Cell.class);
            Cell labelHeader = mock(Cell.class);
            Cell redHeader = mock(Cell.class);
            Cell greenHeader = mock(Cell.class);
            Cell blueHeader = mock(Cell.class);
            Cell valueCell = mock(Cell.class);
            Cell labelCell = mock(Cell.class);
            Cell redCell = mock(Cell.class);
            Cell greenCell = mock(Cell.class);
            Cell blueCell = mock(Cell.class);

            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenReturn(mockWorkbook);
            when(mockWorkbook.getSheet(XlsLegendParser.GLOBAL_LEGEND_SHEET)).thenReturn(mockSheet);

            when(mockSheet.findCell(XlsLegendParser.VALUE)).thenReturn(valueHeader);
            when(valueHeader.getColumn()).thenReturn(0);
            when(mockSheet.findCell(XlsLegendParser.LABEL)).thenReturn(labelHeader);
            when(labelHeader.getColumn()).thenReturn(1);
            when(mockSheet.findCell(XlsLegendParser.RED)).thenReturn(redHeader);
            when(redHeader.getColumn()).thenReturn(2);
            when(mockSheet.findCell(XlsLegendParser.GREEN)).thenReturn(greenHeader);
            when(greenHeader.getColumn()).thenReturn(3);
            when(mockSheet.findCell(XlsLegendParser.BLUE)).thenReturn(blueHeader);
            when(blueHeader.getColumn()).thenReturn(4);

            when(mockSheet.getRows()).thenReturn(3);
            // Row 1 data
            when(mockSheet.getCell(0, 1)).thenReturn(valueCell);
            when(valueCell.getContents()).thenReturn("201");
            when(mockSheet.getCell(1, 1)).thenReturn(labelCell);
            when(labelCell.getContents()).thenReturn("Water Body");
            when(mockSheet.getCell(2, 1)).thenReturn(redCell);
            when(redCell.getContents()).thenReturn("0");
            when(mockSheet.getCell(3, 1)).thenReturn(greenCell);
            when(greenCell.getContents()).thenReturn("0");
            when(mockSheet.getCell(4, 1)).thenReturn(blueCell);
            when(blueCell.getContents()).thenReturn("255");
            // Row 2 data
            Cell valueCell2 = mock(Cell.class);
            Cell labelCell2 = mock(Cell.class);
            Cell redCell2 = mock(Cell.class);
            Cell greenCell2 = mock(Cell.class);
            Cell blueCell2 = mock(Cell.class);
            when(mockSheet.getCell(0, 2)).thenReturn(valueCell2);
            when(valueCell2.getContents()).thenReturn("202");
            when(mockSheet.getCell(1, 2)).thenReturn(labelCell2);
            when(labelCell2.getContents()).thenReturn("Urban Area");
            when(mockSheet.getCell(2, 2)).thenReturn(redCell2);
            when(redCell2.getContents()).thenReturn("128");
            when(mockSheet.getCell(3, 2)).thenReturn(greenCell2);
            when(greenCell2.getContents()).thenReturn("128");
            when(mockSheet.getCell(4, 2)).thenReturn(blueCell2);
            when(blueCell2.getContents()).thenReturn("128");

            LegendClass[] result = parser.parse(mockIs, false);

            assertEquals(2, result.length);
            assertEquals(201, result[0].getValue());
            assertEquals("Water Body", result[0].getDescr());
            assertEquals(new Color(0, 0, 255), result[0].getColor());
            assertEquals(202, result[1].getValue());
            assertEquals("Urban Area", result[1].getDescr());
            assertEquals(new Color(128, 128, 128), result[1].getColor());
            verify(mockWorkbook, times(1)).close();
        }
    }

    @Test
    public void testEmptyValueCellSkipsRow() throws Exception {
        InputStream mockIs = new ByteArrayInputStream(new byte[0]);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            Workbook mockWorkbook = mock(Workbook.class);
            Sheet mockSheet = mock(Sheet.class);
            Cell valueHeader = mock(Cell.class);
            Cell labelHeader = mock(Cell.class);
            Cell redHeader = mock(Cell.class);
            Cell greenHeader = mock(Cell.class);
            Cell blueHeader = mock(Cell.class);

            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenReturn(mockWorkbook);
            when(mockWorkbook.getSheet(XlsLegendParser.GLOBAL_LEGEND_SHEET)).thenReturn(mockSheet);

            when(mockSheet.findCell(XlsLegendParser.VALUE)).thenReturn(valueHeader);
            when(valueHeader.getColumn()).thenReturn(0);
            when(mockSheet.findCell(XlsLegendParser.LABEL)).thenReturn(labelHeader);
            when(labelHeader.getColumn()).thenReturn(1);
            when(mockSheet.findCell(XlsLegendParser.RED)).thenReturn(redHeader);
            when(redHeader.getColumn()).thenReturn(2);
            when(mockSheet.findCell(XlsLegendParser.GREEN)).thenReturn(greenHeader);
            when(greenHeader.getColumn()).thenReturn(3);
            when(mockSheet.findCell(XlsLegendParser.BLUE)).thenReturn(blueHeader);
            when(blueHeader.getColumn()).thenReturn(4);

            when(mockSheet.getRows()).thenReturn(3);
            // Row 1: empty value
            Cell emptyValueCell = mock(Cell.class);
            when(mockSheet.getCell(0, 1)).thenReturn(emptyValueCell);
            when(emptyValueCell.getContents()).thenReturn("");
            // Row 2: valid value
            Cell validValueCell = mock(Cell.class);
            Cell labelCell = mock(Cell.class);
            Cell redCell = mock(Cell.class);
            Cell greenCell = mock(Cell.class);
            Cell blueCell = mock(Cell.class);
            when(mockSheet.getCell(0, 2)).thenReturn(validValueCell);
            when(validValueCell.getContents()).thenReturn("301");
            when(mockSheet.getCell(1, 2)).thenReturn(labelCell);
            when(labelCell.getContents()).thenReturn("Mountain");
            when(mockSheet.getCell(2, 2)).thenReturn(redCell);
            when(redCell.getContents()).thenReturn("139");
            when(mockSheet.getCell(3, 2)).thenReturn(greenCell);
            when(greenCell.getContents()).thenReturn("69");
            when(mockSheet.getCell(4, 2)).thenReturn(blueCell);
            when(blueCell.getContents()).thenReturn("19");

            LegendClass[] result = parser.parse(mockIs, false);
            assertEquals(1, result.length);
            assertEquals(301, result[0].getValue());
            assertEquals("Mountain", result[0].getDescr());
        }
    }

    @Test
    public void testOnlyHeaderRowReturnsEmptyArray() throws Exception {
        InputStream mockIs = new ByteArrayInputStream(new byte[0]);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            Workbook mockWorkbook = mock(Workbook.class);
            Sheet mockSheet = mock(Sheet.class);
            Cell valueHeader = mock(Cell.class);
            Cell labelHeader = mock(Cell.class);
            Cell redHeader = mock(Cell.class);
            Cell greenHeader = mock(Cell.class);
            Cell blueHeader = mock(Cell.class);

            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenReturn(mockWorkbook);
            when(mockWorkbook.getSheet(XlsLegendParser.GLOBAL_LEGEND_SHEET)).thenReturn(mockSheet);

            when(mockSheet.findCell(XlsLegendParser.VALUE)).thenReturn(valueHeader);
            when(valueHeader.getColumn()).thenReturn(0);
            when(mockSheet.findCell(XlsLegendParser.LABEL)).thenReturn(labelHeader);
            when(labelHeader.getColumn()).thenReturn(1);
            when(mockSheet.findCell(XlsLegendParser.RED)).thenReturn(redHeader);
            when(redHeader.getColumn()).thenReturn(2);
            when(mockSheet.findCell(XlsLegendParser.GREEN)).thenReturn(greenHeader);
            when(greenHeader.getColumn()).thenReturn(3);
            when(mockSheet.findCell(XlsLegendParser.BLUE)).thenReturn(blueHeader);
            when(blueHeader.getColumn()).thenReturn(4);

            when(mockSheet.getRows()).thenReturn(1);

            LegendClass[] result = parser.parse(mockIs, false);
            assertEquals(0, result.length);
        }
    }

    @Test
    public void testBiffExceptionReturnsEmptyArray() throws Exception {
        InputStream mockIs = new ByteArrayInputStream(new byte[0]);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenThrow(new BiffException("Invalid XLS format"));
            LegendClass[] result = parser.parse(mockIs, false);
            assertEquals(0, result.length);
        }
    }

    @Test
    public void testIOExceptionReturnsEmptyArray() throws Exception {
        InputStream mockIs = new ByteArrayInputStream(new byte[0]);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenThrow(new IOException("File read error"));
            LegendClass[] result = parser.parse(mockIs, false);
            assertEquals(0, result.length);
        }
    }
}