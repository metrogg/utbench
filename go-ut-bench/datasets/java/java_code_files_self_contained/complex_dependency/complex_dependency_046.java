// Converted Java method
import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class ObjectInspector {

    /**
     * Inspects an object and returns a detailed map of its structure and values.
     * Handles primitive types, strings, arrays, collections, and nested objects.
     * 
     * @param obj The object to inspect
     * @param maxDepth Maximum recursion depth for nested objects
     * @return Map containing field names, types, and values
     * @throws IllegalAccessException If field access is denied
     */
    public static Map<String, Object> inspectObject(Object obj, int maxDepth) throws IllegalAccessException {
        Map<String, Object> result = new HashMap<>();
        if (obj == null || maxDepth < 0) {
            result.put("value", null);
            return result;
        }

        Class<?> clazz = obj.getClass();
        
        // Handle primitive types and strings
        if (isSimpleType(obj)) {
            result.put("type", clazz.getSimpleName());
            result.put("value", obj);
            return result;
        }

        // Handle arrays
        if (clazz.isArray()) {
            List<Object> arrayValues = new ArrayList<>();
            int length = java.lang.reflect.Array.getLength(obj);
            for (int i = 0; i < length; i++) {
                Object element = java.lang.reflect.Array.get(obj, i);
                arrayValues.add(inspectObject(element, maxDepth - 1));
            }
            result.put("type", clazz.getComponentType().getSimpleName() + "[]");
            result.put("value", arrayValues);
            return result;
        }

        // Handle complex objects
        Map<String, Object> fieldsMap = new HashMap<>();
        List<Field> fields = getAllFields(clazz);
        
        for (Field field : fields) {
            if (Modifier.isStatic(field.getModifiers())) {
                continue; // Skip static fields
            }
            
            field.setAccessible(true);
            Object fieldValue = field.get(obj);
            String fieldName = field.getName();
            
            if (isSimpleType(fieldValue)) {
                fieldsMap.put(fieldName, createSimpleFieldMap(field, fieldValue));
            } else {
                fieldsMap.put(fieldName, inspectObject(fieldValue, maxDepth - 1));
            }
        }

        result.put("type", clazz.getSimpleName());
        result.put("fields", fieldsMap);
        return result;
    }

    private static boolean isSimpleType(Object obj) {
        return obj == null || 
               obj instanceof String || 
               obj instanceof Number || 
               obj instanceof Boolean || 
               obj instanceof Character;
    }

    private static Map<String, Object> createSimpleFieldMap(Field field, Object value) {
        Map<String, Object> fieldMap = new HashMap<>();
        fieldMap.put("type", field.getType().getSimpleName());
        fieldMap.put("value", value);
        return fieldMap;
    }

    private static List<Field> getAllFields(Class<?> clazz) {
        List<Field> fields = new ArrayList<>();
        while (clazz != null) {
            fields.addAll(Arrays.asList(clazz.getDeclaredFields()));
            clazz = clazz.getSuperclass();
        }
        return fields;
    }
}
