import org.junit.Test;
import static org.junit.Assert.*;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

public class NameParserTest {
    private final NameParser nameParser = new NameParser();

    @Test(expected = NullPointerException.class)
    public void testParse_NullInput_ThrowsNullPointerException() {
        nameParser.parse(null);
    }

    @Test
    public void testParse_EmptyString_ReturnsEmptyCollection() {
        Collection<NameToken> result = nameParser.parse("");
        assertTrue(result.isEmpty());
    }

    @Test
    public void testParse_OnlySlashesAndWhitespace_ReturnsEmptyCollection() {
        Collection<NameToken> result = nameParser.parse("  ///  //  ");
        assertTrue(result.isEmpty());
    }

    @Test
    public void testParse_AllFixedTokens_ReturnsCorrectFixedTokens() {
        Collection<NameToken> result = nameParser.parse("user/profile/photo");
        assertEquals(3, result.size());
        List<NameToken> tokenList = new ArrayList<>(result);
        
        assertTrue(tokenList.get(0) instanceof FixedToken);
        assertEquals("user", ((FixedToken) tokenList.get(0)).getValue());
        assertTrue(tokenList.get(1) instanceof FixedToken);
        assertEquals("profile", ((FixedToken) tokenList.get(1)).getValue());
        assertTrue(tokenList.get(2) instanceof FixedToken);
        assertEquals("photo", ((FixedToken) tokenList.get(2)).getValue());
    }

    @Test
    public void testParse_WithNumericRange_ReturnsMixedTokenTypes() {
        Collection<NameToken> result = nameParser.parse("item/[1-100]/detail");
        assertEquals(3, result.size());
        List<NameToken> tokenList = new ArrayList<>(result);
        
        assertTrue(tokenList.get(0) instanceof FixedToken);
        assertEquals("item", ((FixedToken) tokenList.get(0)).getValue());
        assertTrue(tokenList.get(1) instanceof NumericRange);
        assertEquals("[1-100]", ((NumericRange) tokenList.get(1)).getValue());
        assertTrue(tokenList.get(2) instanceof FixedToken);
        assertEquals("detail", ((FixedToken) tokenList.get(2)).getValue());
    }

    @Test
    public void testParse_PartsWithWhitespace_TrimsPartsCorrectly() {
        Collection<NameToken> result = nameParser.parse("  order  /  [20-30]  /  invoice  ");
        assertEquals(3, result.size());
        List<NameToken> tokenList = new ArrayList<>(result);
        
        assertEquals("order", ((FixedToken) tokenList.get(0)).getValue());
        assertEquals("[20-30]", ((NumericRange) tokenList.get(1)).getValue());
        assertEquals("invoice", ((FixedToken) tokenList.get(2)).getValue());
    }

    @Test
    public void testParse_LeadingAndTrailingSlashes_IgnoresEmptyParts() {
        Collection<NameToken> result = nameParser.parse("/prefix/[5-9]/suffix/");
        assertEquals(3, result.size());
        List<NameToken> tokenList = new ArrayList<>(result);
        
        assertEquals("prefix", ((FixedToken) tokenList.get(0)).getValue());
        assertEquals("[5-9]", ((NumericRange) tokenList.get(1)).getValue());
        assertEquals("suffix", ((FixedToken) tokenList.get(2)).getValue());
    }

    @Test
    public void testParse_PartStartsWithBracket_NoClosingBracket_CreatesNumericRange() {
        Collection<NameToken> result = nameParser.parse("[invalidRange/test");
        assertEquals(2, result.size());
        List<NameToken> tokenList = new ArrayList<>(result);
        
        assertTrue(tokenList.get(0) instanceof NumericRange);
        assertEquals("[invalidRange", ((NumericRange) tokenList.get(0)).getValue());
        assertTrue(tokenList.get(1) instanceof FixedToken);
        assertEquals("test", ((FixedToken) tokenList.get(1)).getValue());
    }

    @Test
    public void testParse_SingleFixedToken_ReturnsSingleToken() {
        Collection<NameToken> result = nameParser.parse("helloWorld");
        assertEquals(1, result.size());
        NameToken token = result.iterator().next();
        assertTrue(token instanceof FixedToken);
        assertEquals("helloWorld", ((FixedToken) token).getValue());
    }

    @Test
    public void testParse_ConsecutiveSlashes_IgnoresEmptyParts() {
        Collection<NameToken> result = nameParser.parse("a//b///c");
        assertEquals(3, result.size());
        List<NameToken> tokenList = new ArrayList<>(result);
        assertEquals("a", ((FixedToken) tokenList.get(0)).getValue());
        assertEquals("b", ((FixedToken) tokenList.get(1)).getValue());
        assertEquals("c", ((FixedToken) tokenList.get(2)).getValue());
    }
}