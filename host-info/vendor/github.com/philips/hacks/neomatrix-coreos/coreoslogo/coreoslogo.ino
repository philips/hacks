#include <Adafruit_NeoPixel.h>
#include <Adafruit_GFX.h>
#include <Adafruit_NeoMatrix.h>

#define PIN 1

Adafruit_NeoMatrix matrix = Adafruit_NeoMatrix(8, 8, PIN,
NEO_MATRIX_TOP     + NEO_MATRIX_RIGHT +
NEO_MATRIX_COLUMNS + NEO_MATRIX_PROGRESSIVE,
NEO_GRB            + NEO_KHZ800);

const uint16_t colors[] = {
  matrix.Color(255, 0, 0), matrix.Color(0, 255, 0), matrix.Color(0, 0, 255) };

void setup() {
  matrix.begin();
  matrix.setTextWrap(false);
  matrix.setBrightness(13);
  matrix.setTextColor(colors[0]);
}

int x    = matrix.width();
int pass = 0;

void loop() {
  matrix.fillScreen(0);
  matrix.setCursor(x, 0);
matrix.drawPixel(0, 0, matrix.Color(0xe2, 0xed, 0xf7));
matrix.drawPixel(1, 0, matrix.Color(0x97, 0xbf, 0xe0));
matrix.drawPixel(2, 0, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(3, 0, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(4, 0, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(5, 0, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(6, 0, matrix.Color(0x97, 0xbf, 0xe0));
matrix.drawPixel(7, 0, matrix.Color(0xe2, 0xed, 0xf7));
matrix.drawPixel(0, 1, matrix.Color(0x97, 0xbf, 0xe0));
matrix.drawPixel(1, 1, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(2, 1, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(3, 1, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(4, 1, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(5, 1, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(6, 1, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(7, 1, matrix.Color(0x97, 0xbf, 0xe0));
matrix.drawPixel(0, 2, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(1, 2, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(2, 2, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(3, 2, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(4, 2, matrix.Color(0xff, 0xff, 0xff));
matrix.drawPixel(5, 2, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(6, 2, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(7, 2, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(0, 3, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(1, 3, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(2, 3, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(3, 3, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(4, 3, matrix.Color(0xff, 0xff, 0xff));
matrix.drawPixel(5, 3, matrix.Color(0xff, 0xff, 0xff));
matrix.drawPixel(6, 3, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(7, 3, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(0, 4, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(1, 4, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(2, 4, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(3, 4, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(4, 4, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(5, 4, matrix.Color(0xe9, 0x46, 0x59));
matrix.drawPixel(6, 4, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(7, 4, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(0, 5, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(1, 5, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(2, 5, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(3, 5, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(4, 5, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(5, 5, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(6, 5, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(7, 5, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(0, 6, matrix.Color(0x97, 0xbf, 0xe0));
matrix.drawPixel(1, 6, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(2, 6, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(3, 6, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(4, 6, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(5, 6, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(6, 6, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(7, 6, matrix.Color(0x97, 0xbf, 0xe0));
matrix.drawPixel(0, 7, matrix.Color(0xe2, 0xed, 0xf7));
matrix.drawPixel(1, 7, matrix.Color(0x97, 0xbf, 0xe0));
matrix.drawPixel(2, 7, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(3, 7, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(4, 7, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(5, 7, matrix.Color(0x43, 0x90, 0xcf));
matrix.drawPixel(6, 7, matrix.Color(0x97, 0xbf, 0xe0));
matrix.drawPixel(7, 7, matrix.Color(0xe2, 0xed, 0xf7));


  matrix.show();
  delay(100);
}
