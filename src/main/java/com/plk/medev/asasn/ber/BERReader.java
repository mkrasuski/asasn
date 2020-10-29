package com.plk.medev.asasn.ber;

public class BERReader {

    private int ptr;
    private byte[] buffer;

    int next() { return buffer[ptr++]; }

    public Tag readTag() {
        int b = next();
        int cls = (b & Tag.CLASS_MASK);
        long value = (b & Tag.VALUE_MASK);
        if (value == 31) {
            value = b128();
        }

        return new Tag(new Tag.TagClass(cls), value);
    }

    protected long b128() {
        long value = 0;
        for (int b = next(); 0 != (b & 0x80) ; b = next()) {
            value <<= 7;
            value |= (b & 0x7f);
        }
        return value;
    }
}
