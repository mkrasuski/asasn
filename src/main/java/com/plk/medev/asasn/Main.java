package com.plk.medev.asasn;

import com.plk.medev.asasn.ber.BERReader;
import com.plk.medev.asasn.ber.Tag;

import java.io.*;

public class Main {

    static void printTLV(StringBuilder out, BERReader br, String pfx) {
        while (br.available()) {
            Tag tag = br.readTag();
            int len = br.readLen();
            String tagStr = tag.toString();
            out.append(pfx).append(tagStr).append(':').append(len).append(':');
            if (tag.constructed) {
                out.append('\n');
                printTLV(out, br.valueReader(len), pfx + tagStr + ".");
            } else
                out.append(br.readHexBuffer(len)).append('\n');
        }
    }


    public static void main(String[] args) throws IOException {

        String indent = "";
        InputStream is = new BufferedInputStream(
                new FileInputStream("/home/mkr/EPG_FILE_SAMPLE"), 1024*1024);

        byte[] buffer = new byte[50*1024*1024];

        int size = is.read(buffer);

        BERReader br = new BERReader(buffer, 0, size);
        StringBuilder out = new StringBuilder(100 * 1024 *1024);

        try {
            printTLV(out, br, "");
        }
        catch (Exception e) {
            throw new RuntimeException(e);
        }
        System.out.println(out);
    }
}
