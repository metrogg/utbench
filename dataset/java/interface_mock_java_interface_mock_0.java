class XlsLegendParser implements Parser {

    @Override
    public LegendClass[] parse(InputStream inputStream, boolean isRegional) {
        Workbook workbook = null;
        LegendClass[] classes = new LegendClass[0];
        try {
            workbook = Workbook.getWorkbook(inputStream);
            Sheet sheet;
            if( isRegional ) {
                sheet = workbook.getSheet(REGIONAL_LEGEND_SHEET);
                if(sheet == null) {
                    throw new IllegalArgumentException("Given xls-legend is not regional.");
                }
            } else {
                sheet = workbook.getSheet(GLOBAL_LEGEND_SHEET);
            }

            final int valueCol = sheet.findCell(VALUE).getColumn();
            final int labelCol = sheet.findCell(LABEL).getColumn();
            final int redCol = sheet.findCell(RED).getColumn();
            final int blueCol = sheet.findCell(BLUE).getColumn();
            final int greenCol = sheet.findCell(GREEN).getColumn();

            classes = new LegendClass[ sheet.getRows() - 1 ];

            for (int i = 1; i < sheet.getRows(); i++) {
                final Cell valueCell = sheet.getCell(valueCol, i);
                if(StringUtils.isNullOrEmpty(valueCell.getContents())) {
                    continue;
                }
                final Cell labelCell = sheet.getCell(labelCol, i);
                final Cell redCell = sheet.getCell(redCol, i);
                final Cell greenCell = sheet.getCell(greenCol, i);
                final Cell blueCell = sheet.getCell(blueCol, i);

                final int value = Integer.parseInt(valueCell.getContents());
                final String name = "Class_" + i;
                final String descr = labelCell.getContents().trim();
                final Color color = new Color(Integer.parseInt(redCell.getContents()),
                                              Integer.parseInt(greenCell.getContents()),
                                              Integer.parseInt(blueCell.getContents()));

                final LegendClass legendClass = new LegendClass(value, descr, name, color);
                classes[ i - 1 ] = legendClass;
            }

        } catch (BiffException e) {
            Debug.trace(e);
        } catch (IOException e) {
            Debug.trace(e);
        } finally {
            if( workbook != null ) {
                workbook.close();
            }
        }
        return classes;
    }

}

class XlsLegendParserTest {

    private XlsLegendParser parser;

    @Test(expected = IllegalArgumentException.class)
    public void testGlobalLegendParsingRegionalSheet() {
