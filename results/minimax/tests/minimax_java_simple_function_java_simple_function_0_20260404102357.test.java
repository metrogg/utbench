import org.junit.Test;
import static org.junit.Assert.*;
import java.util.Collection;
import java.util.ArrayList;

class NameParserTest {

    @Test
    public void testParse() {
        NameParser parser = new NameParser();
        
        // Test basic parsing
        Collection<NameToken> result = parser.parse("hello/world");
        assertEquals(2, result.size());
        
        // Test single part
        Collection<NameToken> single = parser.parse("hello");
        assertEquals(1, single.size());
        
        // Test empty string
        Collection<NameToken> empty = parser.parse("");
        assertEquals(0, empty.size());
        
        // Test parts with spaces
        Collection<NameToken> spaces = parser.parse("  hello  /  world  ");
        assertEquals(2, spaces.size());
        
        // Test numeric range (starts with [)
        Collection<NameToken> numeric = parser.parse("[0-5]");
        assertEquals(1, numeric.size());
        
        // Test mixed parts
        Collection<NameToken> mixed = parser.parse("name/[0-9]");
        assertEquals(2, mixed.size());
        
        // Test multiple empty parts (should be skipped)
        Collection<NameToken> multipleEmpty = parser.parse("a//b");
        assertEquals(2, multipleEmpty.size());
    }