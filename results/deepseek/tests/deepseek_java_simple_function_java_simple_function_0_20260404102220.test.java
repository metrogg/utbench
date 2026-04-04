import org.junit.Test;
import org.junit.Before;
import static org.junit.Assert.*;
import java.util.Collection;
import java.util.ArrayList;

class NameParserTest {

    private NameParser parser;

    @Before
    public void setUp() {
        parser = new NameParser();
    }

    @Test
    public void testParse_EmptyString() {
        Collection<NameToken> result = parser.parse("");
        assertNotNull(result);
        assertEquals(0, result.size());
    }

    @Test
    public void testParse_SingleFixedToken() {
        Collection<NameToken> result = parser.parse("test");
        assertEquals(1, result.size());
        NameToken token = result.iterator().next();
        assertTrue(token instanceof FixedToken);
        assertEquals("test", ((FixedToken) token).getValue());
    }

    @Test
    public void testParse_SingleFixedTokenWithSpaces() {
        Collection<NameToken> result = parser.parse("  test  ");
        assertEquals(1, result.size());
        NameToken token = result.iterator().next();
        assertTrue(token instanceof FixedToken);
        assertEquals("test", ((FixedToken) token).getValue());
    }

    @Test
    public void testParse_SingleNumericRange() {
        Collection<NameToken> result = parser.parse("[1-10]");
        assertEquals(1, result.size());
        NameToken token = result.iterator().next();
        assertTrue(token instanceof NumericRange);
    }

    @Test
    public void testParse_MultipleFixedTokens() {
        Collection<NameToken> result = parser.parse("folder/file/name");
        assertEquals(3, result.size());
        
        ArrayList<NameToken> tokens = new ArrayList<>(result);
        assertTrue(tokens.get(0) instanceof FixedToken);
        assertEquals("folder", ((FixedToken) tokens.get(0)).getValue());
        
        assertTrue(tokens.get(1) instanceof FixedToken);
        assertEquals("file", ((FixedToken) tokens.get(1)).getValue());
        
        assertTrue(tokens.get(2) instanceof FixedToken);
        assertEquals("name", ((FixedToken) tokens.get(2)).getValue());
    }

    @Test
    public void testParse_MixedTokens() {
        Collection<NameToken> result = parser.parse("folder/[1-5]/file");
        assertEquals(3, result.size());
        
        ArrayList<NameToken> tokens = new ArrayList<>(result);
        assertTrue(tokens.get(0) instanceof FixedToken);
        assertEquals("folder", ((FixedToken) tokens.get(0)).getValue());
        
        assertTrue(tokens.get(1) instanceof NumericRange);
        
        assertTrue(tokens.get(2) instanceof FixedToken);
        assertEquals("file", ((FixedToken) tokens.get(2)).getValue());
    }

    @Test
    public void testParse_WithEmptyParts() {
        Collection<NameToken> result = parser.parse("folder//file");
        assertEquals(2, result.size());
        
        ArrayList<NameToken> tokens = new ArrayList<>(result);
        assertTrue(tokens.get(0) instanceof FixedToken);
        assertEquals("folder", ((FixedToken) tokens.get(0)).getValue());
        
        assertTrue(tokens.get(1) instanceof FixedToken);
        assertEquals("file", ((FixedToken) tokens.get(1)).getValue());
    }

    @Test
    public void testParse_WithLeadingAndTrailingSlashes() {
        Collection<NameToken> result = parser.parse("/folder/file/");
        assertEquals(2, result.size());
        
        ArrayList<NameToken> tokens = new ArrayList<>(result);
        assertTrue(tokens.get(0) instanceof FixedToken);
        assertEquals("folder", ((FixedToken) tokens.get(0)).getValue());
        
        assertTrue(tokens.get(1) instanceof FixedToken);
        assertEquals("file", ((FixedToken) tokens.get(1)).getValue());
    }

    @Test
    public void testParse_OnlySlashes() {
        Collection<NameToken> result = parser.parse("///");
        assertNotNull(result);
        assertEquals(0, result.size());
    }

    @Test
    public void testParse_WhitespaceOnlyParts() {
        Collection<NameToken> result = parser.parse("  /  /  ");
        assertNotNull(result);
        assertEquals(0, result.size());
    }

    @Test
    public void testParse_NumericRangeWithSpaces() {
        Collection<NameToken> result = parser.parse("  [1-10]  ");
        assertEquals(1, result.size());
        NameToken token = result.iterator().next();
        assertTrue(token instanceof NumericRange);
    }

    @Test
    public void testParse_ComplexMixedTokens() {
        Collection<NameToken> result = parser.parse("  dir/[1-5]/sub/[10-20]/file.txt  ");
        assertEquals(4, result.size());
        
        ArrayList<NameToken> tokens = new ArrayList<>(result);
        assertTrue(tokens.get(0) instanceof FixedToken);
        assertEquals("dir", ((FixedToken) tokens.get(0)).getValue());
        
        assertTrue(tokens.get(1) instanceof NumericRange);
        
        assertTrue(tokens.get(2) instanceof FixedToken);
        assertEquals("sub", ((FixedToken) tokens.get(2)).getValue());
        
        assertTrue(tokens.get(3) instanceof NumericRange);
    }

    @Test
    public void testParse_NullInput() {
        try {
            parser.parse(null);
            fail("Expected NullPointerException");
        } catch (NullPointerException e) {
            // Expected exception
        }
    }

    @Test
    public void testParse_SingleCharacterTokens() {
        Collection<NameToken> result = parser.parse("a/b/c");
        assertEquals(3, result.size());
        
        ArrayList<NameToken> tokens = new ArrayList<>(result);
        for (int i = 0; i < 3; i++) {
            assertTrue(tokens.get(i) instanceof FixedToken);
            assertEquals(String.valueOf((char)('a' + i)), ((FixedToken) tokens.get(i)).getValue());
        }
    }

    @Test
    public void testParse_LongTokenWithSpecialCharacters() {
        Collection<NameToken> result = parser.parse("my-folder_name/file.name-v1.0");
        assertEquals(2, result.size());
        
        ArrayList<NameToken> tokens = new ArrayList<>(result);
        assertTrue(tokens.get(0) instanceof FixedToken);
        assertEquals("my-folder_name", ((FixedToken) tokens.get(0)).getValue());
        
        assertTrue(tokens.get(1) instanceof FixedToken);
        assertEquals("file.name-v1.0", ((FixedToken) tokens.get(1)).getValue());
    }
}