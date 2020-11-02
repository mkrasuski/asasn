package com.plk.medev.asasn.ber;

public class BERReader {

    private static final int MASK_7BITS = 0b01111111;
    private static final int MASK_8TH_BIT = 0b10000000;
    private static final int MASK_LO_NIBBLE = 0b00001111;
    private static final int MASK_HI_NIBBLE = 0b11110000;
    private static final int MASK_BYTE = MASK_HI_NIBBLE | MASK_LO_NIBBLE;

    private int ptr;
    private byte[] buffer;
    private int size;

    public BERReader(byte[] buf, int off, int len) {
        buffer = buf;
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
        return (ptr < size) ? (buffer[ptr++] & MASK_BYTE) : -1;
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
            value |= (b & MASK_7BITS);
            if (0 == (b & MASK_8TH_BIT)) break;
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

    public static void printTLV(StringBuilder out, BERReader br, String path, String pfx) {

        while (br.available()) {
            Tag tag = br.readTag();
            int len = br.readLen();

            String fieldPath = path + tag.toString();
            final EPGSchema.ASNField field = EPGSchema.findField(fieldPath);

            out.append(pfx)
                    .append(field == null ? fieldPath : field.field)
                    .append('(')
                    .append(len)
                    .append("):  ");

            if (tag.constructed) {

                out.append('\n');
                printTLV(out, br.valueReader(len), fieldPath + ".", pfx + "  ");

            } else {
                EPGSchema.Stringer stringer = (field != null) ?
                        field.stringer
                        : EPGSchema.valSTRING;
                stringer.stringify(br.buffer, br.ptr, len, out);
                out.append('\n');
                br.skip(len);
            }
        }
    }

    public boolean available() {
        return ptr < size;
    }
}
