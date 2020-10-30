package com.plk.medev.asasn.ber;

public class BERReader {

    private static final String HEX_DIGITS = "0123456789ABCDEF";

    private int offset;
    private int ptr;
    private byte[] buffer;
    private int size;

    public BERReader(byte[] buf, int off, int len) {
        buffer = buf;
        offset = off;
        ptr = off;
        size = off + len;
    }

    public BERReader(byte[] buf) {
        this(buf, 0, buf.length);
    }

    BERReader(BERReader br, int len) {
        this(br.buffer, br.ptr, len);
    }

    private int next() {
        return (ptr < size) ? (buffer[ptr++] & 0xFF) : -1;
    }

    public Tag readTag() {
        int b = next();
        long value = (b & Tag.VALUE_MASK);
        if (value == 31) {
            value = readBase128();
        }

        return new Tag(
                (b & Tag.CLASS_MASK),
                0 != (b & Tag.CONSTRUCT_MASK),
                value);
    }

    public int readLen() {

        int b = next();

        assert (b != 128); // indefinite length

        if (b < 128)
            return b;
        else {
            b -= 128;
            int len = 0;
            while (b-- > 0) {
                len <<= 8;
                len |= next();
            }
            return len;
        }
    }

    protected long readBase128() {
        long value = 0;
        while (true) {
            int b = next();
            value |= (b & 0x7f);
            if (0 == (b & 0x80)) break;
            value <<= 7;
        }
        return value;
    }

    public void skip(long count) {
        ptr += count;
    }

    public BERReader valueReader(int len) {
        BERReader vr = new BERReader(this, len);
        ptr += len;
        return vr;
    }

    public String readHexBuffer(int len) {
        char[] hex = new char[2*len];
        for (int i = 0, p = 0; i < len; i++) {
            int b = next();
            hex[p++] = HEX_DIGITS.charAt((b & 0xF0) >> 4);
            hex[p++] = HEX_DIGITS.charAt(b & 0x0F);
        }
        return String.valueOf(hex);
    }


    public boolean available() {
        return ptr < size;
    }
}
