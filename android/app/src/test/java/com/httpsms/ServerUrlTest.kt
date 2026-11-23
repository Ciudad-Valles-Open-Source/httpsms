package com.nerus.httpsms

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Test

class ServerUrlTest {
    @Test
    fun normalize_trimsSlashesAndV1() {
        assertEquals("https://api.example.com", ServerUrl.normalize("  https://api.example.com/  "))
        assertEquals("https://api.example.com", ServerUrl.normalize("https://api.example.com/v1"))
        assertEquals("https://api.example.com", ServerUrl.normalize("https://api.example.com/V1/"))
        assertEquals("http://192.168.1.10:8000/sms", ServerUrl.normalize("http://192.168.1.10:8000/sms/v1/"))
        assertEquals("", ServerUrl.normalize(null))
    }

    @Test
    fun validate_acceptsHttpAndHttps() {
        assertEquals(ServerUrl.Validation.VALID, ServerUrl.validate("https://api.example.com"))
        assertEquals(ServerUrl.Validation.VALID, ServerUrl.validate("http://192.168.1.10:8000"))
        assertEquals(ServerUrl.Validation.EMPTY, ServerUrl.validate(""))
        assertEquals(ServerUrl.Validation.INVALID, ServerUrl.validate("api.example.com"))
        assertEquals(ServerUrl.Validation.INVALID, ServerUrl.validate("ftp://api.example.com"))
        assertEquals(ServerUrl.Validation.INVALID, ServerUrl.validate("https://"))
        assertEquals(ServerUrl.Validation.INVALID, ServerUrl.validate("https://api.example.com?x=1"))
    }

    @Test
    fun warnings() {
        assertTrue(ServerUrl.isInsecure("http://192.168.1.10:8000"))
        assertFalse(ServerUrl.isInsecure("https://api.example.com"))
        assertTrue(ServerUrl.isLoopback("http://localhost:8000"))
        assertTrue(ServerUrl.isLoopback("http://127.0.0.1:8000"))
        assertFalse(ServerUrl.isLoopback("http://192.168.1.10:8000"))
    }

    @Test
    fun parseLoginQr_newFormat() {
        val qr = ServerUrl.parseLoginQr(
            "httpsms://login?server=https%3A%2F%2Fapi.example.com%2Fv1%2F&key=pk_abc%2B123%3D"
        )
        assertEquals("pk_abc+123=", qr?.apiKey)
        assertEquals("https://api.example.com", qr?.serverUrl)
    }

    @Test
    fun parseLoginQr_newFormatWithoutServer() {
        val qr = ServerUrl.parseLoginQr("httpsms://login?key=pk_abc")
        assertEquals("pk_abc", qr?.apiKey)
        assertNull(qr?.serverUrl)
    }

    @Test
    fun parseLoginQr_oldFormatBareKey() {
        val qr = ServerUrl.parseLoginQr("  pk_abc123  ")
        assertEquals("pk_abc123", qr?.apiKey)
        assertNull(qr?.serverUrl)
    }

    @Test
    fun parseLoginQr_invalid() {
        assertNull(ServerUrl.parseLoginQr(""))
        assertNull(ServerUrl.parseLoginQr(null))
        assertNull(ServerUrl.parseLoginQr("httpsms://login?server=https%3A%2F%2Fapi.example.com"))
        assertNull(ServerUrl.parseLoginQr("httpsms://other?key=pk_abc"))
    }
}
