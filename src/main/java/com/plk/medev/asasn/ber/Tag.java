package com.plk.medev.asasn.ber;

public class Tag {

    public static final int UNIVERSAL = 0b00000000;
    public static final int APPLICATION = 0b01000000;
    public static final int CONTEXT = 0b10000000;
    public static final int PRIVATE = 0b11000000;
    public static final int CLASS_MASK = 0b11000000;
    public static final int VALUE_MASK = 0b00011111;
    public static final int CONSTRUCT_MASK = 0b00100000;


    public int tagClass;
    public boolean constructed;
    public long tagValue;

    Tag(int cls, boolean cnst, long val) {
        tagClass = cls;
        constructed = cnst;
        tagValue = val;
    }

    @Override
    public String toString() {
        return className[tagClass >> 6] + tagValue;
    }

    static final String[] className = {"U", "A", "C", "P"};
}
