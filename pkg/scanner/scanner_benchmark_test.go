package scanner

import (
	"strings"
	"testing"

	"github.com/mkeesey/craftinginterpreters/pkg/failure"
)

func BenchmarkScanner(b *testing.B) {
	reporter := &failure.Reporter{}
	reader := strings.NewReader(program)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, 0)
		scan := NewScanner(reader, reporter)
		scan.ScanTokens()
	}
}

const program = `
class Zoo {
  init() {
    this.aardvark = 1;
    this.baboon   = 1;
    this.cat      = 1;
    this.donkey   = 1;
    this.elephant = 1;
    this.fox      = 1;
  }
  ant()    { return this.aardvark; }
  banana() { return this.baboon; }
  tuna()   { return this.cat; }
  hay()    { return this.donkey; }
  grass()  { return this.elephant; }
  mouse()  { return this.fox; }
}

var zoo = Zoo();
var sum = 0;
var start = clock();
while (sum < 100000000) {
  sum = sum + zoo.ant()
            + zoo.banana()
            + zoo.tuna()
            + zoo.hay()
            + zoo.grass()
            + zoo.mouse();
}

print clock() - start;
print sum;
`
