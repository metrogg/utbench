import org.junit.Test;
import org.junit.Assert;
import java.util.Collection;
import java.util.Iterator;

public class NameParserTest {

    private final NameParser parser = new NameParser();

    @Test
    public void testParseSingleFixedToken() {
        Collection<NameToken> result = parser.parse("hello");
        Assert.assertEquals(1, result.size());
        Assert.assertTrue(result.iterator().next() instanceof FixedToken);
    }

    @Test
    public void testParseMultipleFixedTokens() {
        Collection<NameToken> result = parser.parse("foo/bar");
        Assert.assertEquals(2, result.size());
        Iterator<NameToken> it = result.iterator();
        Assert.assertTrue(it.next() instanceof FixedToken);
        Assert.assertTrue(it.next() instanceof FixedToken);
    }

    @Test
    public void testParseNumericRange() {
        Collection<NameToken> result = parser.parse("[1-5]");
        Assert.assertEquals(1, result.size());
        Assert.assertTrue(result.iterator().next() instanceof NumericRange);
    }

    @Test
    public void testParseMixedTokens() {
        Collection<NameToken> result = parser.parse("start/[0-9]/end");
        Assert.assertEquals(3, result.size());
        Iterator<NameToken> it = result.iterator();
        Assert.assertTrue(it.next() instanceof FixedToken);
        Assert.assertTrue(it.next() instanceof NumericRange);
        Assert.assertTrue(it.next() instanceof FixedToken);
    }

    @Test
    public void testParseWithWhitespaceAndTrimming() {
        Collection<NameToken> result = parser.parse(" a / [1-3] / b ");
        Assert.assertEquals(3, result.size());
        Iterator<NameToken> it = result.iterator();
        Assert.assertTrue(it.next() instanceof FixedToken);
        Assert.assertTrue(it.next() instanceof NumericRange);
        Assert.assertTrue(it.next() instanceof FixedToken);
    }

    @Test
    public void testParseConsecutiveSlashesSkipsEmptyParts() {
        Collection<NameToken> result = parser.parse("a//b");
        Assert.assertEquals(2, result.size());
        Iterator<NameToken> it = result.iterator();
        Assert.assertTrue(it.next() instanceof FixedToken);
        Assert.assertTrue(it.next() instanceof FixedToken);
    }

    @Test
    public void testParseEmptyString() {
        Collection<NameToken> result = parser.parse("");
        Assert.assertEquals(0, result.size());
    }

    @Test
    public void testParseOnlySlashesAndSpaces() {
        Collection<NameToken> result = parser.parse(" / / ");
        Assert.assertEquals(0, result.size());
    }

    @Test
    public void testParseLeadingAndTrailingSlashes() {
        Collection<NameToken> result = parser.parse("/a/b/");
        Assert.assertEquals(2, result.size());
    }

    @Test(expected = NullPointerException.class)
    public void testParseNullInputThrowsException() {
        parser.parse(null);
    }
}