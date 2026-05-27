#!/usr/bin/env python3
"""Create a nice flashcard app icon with matplotlib"""
import os

try:
    import matplotlib.pyplot as plt
    import matplotlib.patches as patches
    from matplotlib.patches import FancyBboxPatch, Circle
    
    # Create figure
    fig, ax = plt.subplots(1, 1, figsize=(8, 8), dpi=64)
    ax.set_xlim(0, 100)
    ax.set_ylim(0, 100)
    ax.axis('off')
    
    # Create gradient background
    gradient = patches.Rectangle((0, 0), 100, 100, 
                                 facecolor='#2196F3', edgecolor='none')
    ax.add_patch(gradient)
    
    # Draw a stylized card/flashcard
    card_width = 60
    card_height = 50
    card_x = (100 - card_width) / 2
    card_y = 25
    
    # Main card with rounded corners
    card = FancyBboxPatch((card_x, card_y), card_width, card_height,
                          boxstyle="round,pad=2", 
                          facecolor='white', 
                          edgecolor='#1976D2',
                          linewidth=3)
    ax.add_patch(card)
    
    # Accent circle
    circle = Circle((70, 30), 8, facecolor='#FF9800', edgecolor='#F57C00', linewidth=2)
    ax.add_patch(circle)
    
    # Draw a checkmark inside the circle
    ax.text(70, 30, '✓', ha='center', va='center', 
            fontsize=20, color='white', weight='bold')
    
    # Add some text on the card
    ax.text(50, 55, 'Learn', ha='center', va='center',
            fontsize=16, color='#2196F3', weight='bold')
    ax.text(50, 40, 'Remember', ha='center', va='center',
            fontsize=12, color='#666666')
    
    # Save figure as PNG
    os.makedirs('assets', exist_ok=True)
    plt.savefig('assets/icon.png', dpi=64, bbox_inches='tight', 
                facecolor='#2196F3', edgecolor='none', pad_inches=0)
    plt.close()
    print('Nice icon created at assets/icon.png')
    
except ImportError:
    print("matplotlib not available, creating icon with PIL instead...")
    try:
        from PIL import Image, ImageDraw, ImageFont
        
        # Create image
        size = (512, 512)
        img = Image.new('RGBA', size, (33, 150, 243, 255))  # Material Blue
        draw = ImageDraw.Draw(img)
        
        # Draw a card shape
        card_box = (60, 120, 452, 392)
        draw.rounded_rectangle(card_box, radius=20, fill='white', outline=(25, 118, 210), width=8)
        
        # Draw an orange circle with checkmark
        circle_bbox = [370, 60, 430, 120]
        draw.ellipse(circle_bbox, fill=(255, 152, 0), outline=(245, 127, 23), width=3)
        
        # Add text
        try:
            draw.text((256, 230), "Learn", fill=(33, 150, 243), anchor="mm", 
                     font=None)  # Using default font
            draw.text((256, 280), "Remember", fill=(100, 100, 100), anchor="mm",
                     font=None)
        except:
            pass
        
        os.makedirs('assets', exist_ok=True)
        img.save('assets/icon.png')
        print('Nice icon created at assets/icon.png (PIL version)')
        
    except ImportError:
        print("PIL also not available, creating gradient icon with struct/zlib...")
        import struct
        import zlib
        
        def create_gradient_png(width, height):
            """Create PNG with gradient background"""
            png_sig = b'\x89PNG\r\n\x1a\n'
            
            # IHDR chunk
            ihdr_data = struct.pack('>IIBBBBB', width, height, 8, 2, 0, 0, 0)
            ihdr_crc = zlib.crc32(b'IHDR' + ihdr_data) & 0xffffffff
            ihdr = struct.pack('>I', 13) + b'IHDR' + ihdr_data + struct.pack('>I', ihdr_crc)
            
            # IDAT chunk with gradient
            raw_data = b''
            for y in range(height):
                raw_data += bytes([0])  # filter type
                ratio = y / height
                # Gradient from #2196F3 (33, 150, 243) to lighter blue
                r = int(33 + (100 * ratio))
                g = int(150 + (50 * ratio))
                b = int(243 - (30 * ratio))
                for x in range(width):
                    raw_data += bytes([r, g, b])
            
            idat_data = zlib.compress(raw_data, 9)
            idat_crc = zlib.crc32(b'IDAT' + idat_data) & 0xffffffff
            idat = struct.pack('>I', len(idat_data)) + b'IDAT' + idat_data + struct.pack('>I', idat_crc)
            
            # IEND chunk
            iend = struct.pack('>I', 0) + b'IEND' + struct.pack('>I', 0xae426082)
            
            return png_sig + ihdr + idat + iend
        
        icon_data = create_gradient_png(512, 512)
        os.makedirs('assets', exist_ok=True)
        with open('assets/icon.png', 'wb') as f:
            f.write(icon_data)
        print('Nice gradient icon created at assets/icon.png')
