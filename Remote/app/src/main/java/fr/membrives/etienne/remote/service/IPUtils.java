package fr.membrives.etienne.remote.service;

import android.util.Log;

import java.math.BigInteger;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.HashSet;
import java.util.Set;

/**
 * Created by etienne on 23/05/15.
 */
public class IPUtils {
    public static Set<String> getIpsInSubnet(int ip, int netmask) {
        int minIp = ip & netmask;
        int maxIp = ip | ~netmask;

        Set<String> ips = new HashSet();

        for (int i = minIp; i <= maxIp; i++) {
            if ((i & 0xFF) == 0x00 || (i & 0xFF) == 0xFF) {
                continue;
            }
            ips.add(ipToString(i));
        }

        return ips;
    }

    private static String ipToString(int ipAddress) {
        byte[] ipByteArray = BigInteger.valueOf(ipAddress).toByteArray();

        String ipAddressString;
        try {
            ipAddressString = InetAddress.getByAddress(ipByteArray).getHostAddress();
        } catch (UnknownHostException e) {
            throw new RuntimeException("IP address: " + Integer.toString(ipAddress), e);
        }
        return ipAddressString;
    }
}
