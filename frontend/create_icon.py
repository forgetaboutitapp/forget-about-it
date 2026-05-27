#!/usr/bin/env python3
"""Create a simple app icon"""
import struct
import zlib
import os

def create_png(width, height, color_rgb):
    """Create a minimal PNG file with solid color"""
    r, g, b = color_rgb
    
    # PNG signature
    png_sig = b'\x89PNG\r\n\x1a\n'
    
    # IHDR chunk (image header)
    ihdr_data = struct.pack('>IIBBBBB', width, height, 8, 2, 0, 0, 0)
    ihdr_crc = zlib.crc32(b'IHDR' + ihdr_data) & 0xffffffff
    ihdr = struct.pack('>I', 13) + b'IHDR' + ihdr_data + struct.pack('>I', ihdr_crc)
    
    # IDAT chunk (image data)
    # Create raw image data: for each scanline, add filter byte 0 (no filter)
    raw_data = b''
    for y in range(height):
        raw_data += bytes([0])  # filter type
        for x in range(width):
            raw_data += bytes([r, g, b])  # RGB pixel
    
    idat_data = zlib.compress(raw_data, 9)
    idat_crc = zlib.crc32(b'IDAT' + idat_data) & 0xffffffff
    idat = struct.pack('>I', len(idat_data)) + b'IDAT' + idat_data + struct.pack('>I', idat_crc)
    
    # IEND chunk (end marker)
    iend = struct.pack('>I', 0) + b'IEND' + struct.pack('>I', 0xae426082)
    
    return png_sig + ihdr + idat + iend

# Create assets directory if it doesn't exist
os.makedirs('assets', exist_ok=True)

# Create a 512x512 blue icon with white "F" in the center
# For simplicity, we'll just create a solid blue square
# (Adding text would require more complex pixel manipulation)
icon_data = create_png(512, 512, (33, 150, 243))  # Material Blue

# Write to file
with open('assets/icon.png', 'wb') as f:
    f.write(icon_data)

print('Icon created at assets/icon.png')
