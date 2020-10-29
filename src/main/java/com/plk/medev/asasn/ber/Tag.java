package com.plk.medev.asasn.ber;

public class Tag {

    public static final int UNIVERSAL = 0b00000000;
    public static final int APPLICATION = 0b01000000;
    public static final int CONTEXT = 0b10000000;
    public static final int PRIVATE = 0b11000000;
    public static final int CLASS_MASK = 0b11000000;
    public static final int VALUE_MASK = 0b00111111;

    public enum TagClass {

        Universal(UNIVERSAL),
        Application(APPLICATION),
        Context(CONTEXT),
        Private(PRIVATE);

        int value;

        TagClass(int tagClass) { value = tagClass; }
    }

    public TagClass tagClass;
    public long tagValue;

    Tag(int c, long value) {
        tagClass = TagClass.Universal;
        tagValue = value;
    }

    @Override
    public String toString() {
        return "Tag{" +
                "tagClass=" + tagClass +
                ", tagValue=" + tagValue +
                '}';
    }
}
