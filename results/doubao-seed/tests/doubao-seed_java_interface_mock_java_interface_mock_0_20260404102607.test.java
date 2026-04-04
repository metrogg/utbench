import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.MockedStatic;
import org.mockito.junit.MockitoJUnitRunner;
import jxl.Workbook;
import jxl.Sheet;
import jxl.Cell;
import jxl.read.biff.BiffException;
import java.io.InputStream;
import java.io.IOException;
import java.awt.Color;
import static org.junit.Assert.*;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@RunWith(MockitoJUnitRunner.class)
class XlsLegendParserTest {

    private XlsLegendParser parser;

    @Before
    public void setUp() {
        parser = new XlsLegendParser();
    }

    @Test(expected = IllegalArgumentException.class)
    public void testGlobalLegendParsingRegionalSheet() throws Exception {
        InputStream mockInputStream = mock(InputStream.class);
        Workbook mockWorkbook = mock(Workbook.class);
        when(mockWorkbook.getSheet(XlsLegendParser.REGIONAL_LEGEND_SHEET)).thenReturn(null);

        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockInputStream)).thenReturn(mockWorkbook);
            parser.parse(mockInputStream, true);
        }
    }

    @Test
    public void testRegionalLegendParseSuccess() throws Exception {
        InputStream mockIs = mock(InputStream.class);
        Workbook mockWb = mock(Workbook.class);
        Sheet mockSheet = mock(Sheet.class);

        Cell valueHeader = mock(Cell.class);
        Cell labelHeader = mock(Cell.class);
        Cell redHeader = mock(Cell.class);
        Cell greenHeader = mock(Cell.class);
        Cell blueHeader = mock(Cell.class);

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

        Cell value1 = mock(Cell.class);
        when(value1.getContents()).thenReturn("1");
        Cell label1 = mock(Cell.class);
        when(label1.getContents()).thenReturn("First Class ");
        Cell red1 = mock(Cell.class);
        when(red1.getContents()).thenReturn("255");
        Cell green1 = mock(Cell.class);
        when(green1.getContents()).thenReturn("0");
        Cell blue1 = mock(Cell.class);
        when(blue1.getContents()).thenReturn("0");

        when(mockSheet.getCell(0, 1)).thenReturn(value1);
        when(mockSheet.getCell(1, 1)).thenReturn(label1);
        when(mockSheet.getCell(2, 1)).thenReturn(red1);
        when(mockSheet.getCell(3, 1)).thenReturn(green1);
        when(mockSheet.getCell(4, 1)).thenReturn(blue1);

        Cell value2 = mock(Cell.class);
        when(value2.getContents()).thenReturn("2");
        Cell label2 = mock(Cell.class);
        when(label2.getContents()).thenReturn("Second Class");
        Cell red2 = mock(Cell.class);
        when(red2.getContents()).thenReturn("0");
        Cell green2 = mock(Cell.class);
        when(green2.getContents()).thenReturn("255");
        Cell blue2 = mock(Cell.class);
        when(blue2.getContents()).thenReturn("0");

        when(mockSheet.getCell(0, 2)).thenReturn(value2);
        when(mockSheet.getCell(1, 2)).thenReturn(label2);
        when(mockSheet.getCell(2, 2)).thenReturn(red2);
        when(mockSheet.getCell(3, 2)).thenReturn(green2);
        when(mockSheet.getCell(4, 2)).thenReturn(blue2);

        when(mockWb.getSheet(XlsLegendParser.REGIONAL_LEGEND_SHEET)).thenReturn(mockSheet);

        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenReturn(mockWb);
            LegendClass[] result = parser.parse(mockIs, true);

            assertEquals(2, result.length);
            assertEquals(1, result[0].getValue());
            assertEquals("First Class", result[0].getDescr());
            assertEquals("Class_1", result[0].getName());
            assertEquals(new Color(255, 0, 0), result[0].getColor());

            assertEquals(2, result[1].getValue());
            assertEquals("Second Class", result[1].getDescr());
            assertEquals("Class_2", result[1].getName());
            assertEquals(new Color(0, 255, 0), result[1].getColor());
        }
    }

    @Test
    public void testGlobalLegendParseWithEmptyValueRowSkipped() throws Exception {
        InputStream mockIs = mock(InputStream.class);
        Workbook mockWb = mock(Workbook.class);
        Sheet mockSheet = mock(Sheet.class);

        Cell valueHeader = mock(Cell.class);
        when(mockSheet.findCell(XlsLegendParser.VALUE)).thenReturn(valueHeader);
        when(valueHeader.getColumn()).thenReturn(0);
        Cell labelHeader = mock(Cell.class);
        when(mockSheet.findCell(XlsLegendParser.LABEL)).thenReturn(labelHeader);
        when(labelHeader.getColumn()).thenReturn(1);
        Cell redHeader = mock(Cell.class);
        when(mockSheet.findCell(XlsLegendParser.RED)).thenReturn(redHeader);
        when(redHeader.getColumn()).thenReturn(2);
        Cell greenHeader = mock(Cell.class);
        when(mockSheet.findCell(XlsLegendParser.GREEN)).thenReturn(greenHeader);
        when(greenHeader.getColumn()).thenReturn(3);
        Cell blueHeader = mock(Cell.class);
        when(mockSheet.findCell(XlsLegendParser.BLUE)).thenReturn(blueHeader);
        when(blueHeader.getColumn()).thenReturn(4);

        when(mockSheet.getRows()).thenReturn(3);

        Cell value1 = mock(Cell.class);
        when(value1.getContents()).thenReturn("");
        when(mockSheet.getCell(0,1)).thenReturn(value1);

        Cell value2 = mock(Cell.class);
        when(value2.getContents()).thenReturn("3");
        Cell label2 = mock(Cell.class);
        when(label2.getContents()).thenReturn("Valid Class");
        Cell red2 = mock(Cell.class);
        when(red2.getContents()).thenReturn("0");
        Cell green2 = mock(Cell.class);
        when(green2.getContents()).thenReturn("0");
        Cell blue2 = mock(Cell.class);
        when(blue2.getContents()).thenReturn("255");

        when(mockSheet.getCell(0,2)).thenReturn(value2);
        when(mockSheet.getCell(1,2)).thenReturn(label2);
        when(mockSheet.getCell(2,2)).thenReturn(red2);
        when(mockSheet.getCell(3,2)).thenReturn(green2);
        when(mockSheet.getCell(4,2)).thenReturn(blue2);

        when(mockWb.getSheet(XlsLegendParser.GLOBAL_LEGEND_SHEET)).thenReturn(mockSheet);

        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenReturn(mockWb);
            LegendClass[] result = parser.parse(mockIs, false);
            assertEquals(1, result.length);
            assertEquals(3, result[0].getValue());
            assertEquals("Valid Class", result[0].getDescr());
        }
    }

    @Test
    public void testBiffExceptionReturnsEmptyArray() throws Exception {
        InputStream mockIs = mock(InputStream.class);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenThrow(new BiffException("Invalid XLS format"));
            LegendClass[] result = parser.parse(mockIs, false);
            assertEquals(0, result.length);
        }
    }

    @Test
    public void testIOExceptionReturnsEmptyArray() throws Exception {
        InputStream mockIs = mock(InputStream.class);
        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenThrow(new IOException("File read error"));
            LegendClass[] result = parser.parse(mockIs, false);
            assertEquals(0, result.length);
        }
    }

    @Test
    public void testWorkbookClosedAfterProcessing() throws Exception {
        InputStream mockIs = mock(InputStream.class);
        Workbook mockWb = mock(Workbook.class);
        Sheet mockSheet = mock(Sheet.class);
        when(mockWb.getSheet(XlsLegendParser.GLOBAL_LEGEND_SHEET)).thenReturn(mockSheet);

        Cell mockHeader = mock(Cell.class);
        when(mockSheet.findCell(anyString())).thenReturn(mockHeader);
        when(mockHeader.getColumn()).thenReturn(0);
        when(mockSheet.getRows()).thenReturn(1);

        try (MockedStatic<Workbook> mockedWorkbook = mockStatic(Workbook.class)) {
            mockedWorkbook.when(() -> Workbook.getWorkbook(mockIs)).thenReturn(mockWb);
            parser.parse(mockIs, false);
            verify(mockWb, times(1)).close();
        }
    }
}