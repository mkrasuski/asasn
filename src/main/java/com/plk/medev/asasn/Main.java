package com.plk.medev.asasn;

import com.plk.medev.asasn.ber.BERReader;

import java.io.BufferedInputStream;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;

public class Main {

    public static void main(String[] args) throws IOException {

        String indent = "";
        InputStream is = new BufferedInputStream(
                new FileInputStream(args[0]), 1024*1024);

        byte[] buffer = new byte[50*1024*1024];

        int size = is.read(buffer);

        BERReader br = new BERReader(buffer, 0, size);
        StringBuilder out = new StringBuilder(50 * 1024 *1024);

        try {
            BERReader.printTLV(out, br, ".", "");
        }
        catch (Exception e) {
            throw new RuntimeException(e);
        }
        System.out.println(out);
    }
}
