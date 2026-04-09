class NameParser {

    public Collection<NameToken> parse(String name) {
        ArrayList<NameToken> list = new ArrayList();
        String[] parts = name.split("/");        
        for (String part: parts) {
            part = part.trim();
            
            if (part.length() == 0) {
                continue;
            }
            
            if (part.startsWith("[")) {
                list.add(new NumericRange(part));
            } else {
                list.add(new FixedToken(part));
            }
            //FIXME: include text ranges
        }
        return list;
    }

}

class NameParserTest {

    @Test
    public void testParse() {
