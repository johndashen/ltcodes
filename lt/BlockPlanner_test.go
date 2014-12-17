package lt

import "testing"

func TestBlockPlanner(t *testing.T) {

	planner := NewBlockPlanner(571, 166362120)
	seeds := []uint32{166362120, 634813345, 177020911, 1055302029, 1364977754,
		1692838451, 915510748, 1536644533, 980758720, 1049939729,
		464738808, 1622156932, 667094411, 526649093, 877036565,
		1891461182, 1813567941, 1687591223, 886846905, 1912570498,
		473728667, 1321711281, 706125047}
	sourceBlocks := []([]uint32){{98}, {400, 62}, {49, 385}, {421, 541},
					   {336,109,412,410,463,231,319,564,417,305,313,461},
					   {444, 522, 416, 49, 9, 199, 239, 182},
					   {370, 167}, {458, 555}, {236, 557, 326, 25, 418, 154, 230, 346},
					   {84, 195}, {138, 177}, {109, 43, 446, 250},
					   {201, 291, 424, 197, 401, 108, 38, 85, 382, 53, 430, 102,117, 454,
						360, 29, 363, 271, 230, 63, 448, 186, 206, 257, 80, 10, 99, 190,
						224, 474, 338, 351, 376},
					   {262, 239, 265, 91, 527, 268, 550}, {271}, {19, 566}, {553, 78, 160, 152},
					   {240, 385, 542, 394, 465, 539},
					   {380, 345, 290, 31, 273, 79, 416, 108, 288}, {129, 204, 230},
					   {326, 461, 451}, {439, 181}, {127, 144}}

	for i, expSeed := range seeds {
		nextList, nextSeed := planner.NextBlockList() 
		if nextSeed != seeds[i] {
			t.Errorf("BlockList #%d seed = %d, expected %d", i, nextSeed, expSeed)
		} else if len(nextList) != len(sourceBlocks[i]) {
			t.Errorf("BlockList #%d len(list) = %d, expected %d", i, len(nextList), len(sourceBlocks[i]))
		} else {
			for idx, val := range nextList {
				if val != sourceBlocks[i][idx] {
					t.Errorf("BlockList #%d list[%d] = %d, expected %d", i, idx, nextList[idx], sourceBlocks[i][idx])
				}
			}
		}
	}
}
