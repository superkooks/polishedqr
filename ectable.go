package polishedqr

type ecBlock struct {
	count     int
	dataWords int
}

type ecBlocks struct {
	ecWordsPerBlock int
	blocks          []ecBlock
}

// i am dead
var codeWordTable = map[int]map[string]ecBlocks{
	1: {
		"L": ecBlocks{ecWordsPerBlock: 7, blocks: []ecBlock{{count: 1, dataWords: 19}}},
		"M": ecBlocks{ecWordsPerBlock: 10, blocks: []ecBlock{{count: 1, dataWords: 16}}},
		"Q": ecBlocks{ecWordsPerBlock: 13, blocks: []ecBlock{{count: 1, dataWords: 13}}},
		"H": ecBlocks{ecWordsPerBlock: 17, blocks: []ecBlock{{count: 1, dataWords: 9}}},
	},
	2: {
		"L": ecBlocks{ecWordsPerBlock: 10, blocks: []ecBlock{{count: 1, dataWords: 34}}},
		"M": ecBlocks{ecWordsPerBlock: 16, blocks: []ecBlock{{count: 1, dataWords: 28}}},
		"Q": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{{count: 1, dataWords: 22}}},
		"H": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{{count: 1, dataWords: 16}}},
	},
	3: {
		"L": ecBlocks{ecWordsPerBlock: 15, blocks: []ecBlock{{count: 1, dataWords: 55}}},
		"M": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{{count: 1, dataWords: 44}}},
		"Q": ecBlocks{ecWordsPerBlock: 18, blocks: []ecBlock{{count: 2, dataWords: 17}}},
		"H": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{{count: 2, dataWords: 13}}},
	},
	4: {
		"L": ecBlocks{ecWordsPerBlock: 20, blocks: []ecBlock{{count: 1, dataWords: 80}}},
		"M": ecBlocks{ecWordsPerBlock: 18, blocks: []ecBlock{{count: 2, dataWords: 32}}},
		"Q": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{{count: 2, dataWords: 24}}},
		"H": ecBlocks{ecWordsPerBlock: 16, blocks: []ecBlock{{count: 4, dataWords: 9}}},
	},
	5: {
		"L": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{{count: 1, dataWords: 108}}},
		"M": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{{count: 2, dataWords: 43}}},
		"Q": ecBlocks{ecWordsPerBlock: 18, blocks: []ecBlock{
			{count: 2, dataWords: 15},
			{count: 2, dataWords: 16},
		}},
		"H": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{
			{count: 2, dataWords: 11},
			{count: 2, dataWords: 12},
		}},
	},
	6: {
		"L": ecBlocks{ecWordsPerBlock: 18, blocks: []ecBlock{{count: 2, dataWords: 68}}},
		"M": ecBlocks{ecWordsPerBlock: 16, blocks: []ecBlock{{count: 4, dataWords: 27}}},
		"Q": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{{count: 4, dataWords: 19}}},
		"H": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{{count: 4, dataWords: 15}}},
	},
	7: {
		"L": ecBlocks{ecWordsPerBlock: 20, blocks: []ecBlock{{count: 2, dataWords: 78}}},
		"M": ecBlocks{ecWordsPerBlock: 18, blocks: []ecBlock{{count: 4, dataWords: 31}}},
		"Q": ecBlocks{ecWordsPerBlock: 18, blocks: []ecBlock{
			{count: 2, dataWords: 14},
			{count: 4, dataWords: 15},
		}},
		"H": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 4, dataWords: 13},
			{count: 1, dataWords: 14},
		}},
	},
	8: {
		"L": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{{count: 2, dataWords: 97}}},
		"M": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{
			{count: 2, dataWords: 38},
			{count: 2, dataWords: 39},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{
			{count: 4, dataWords: 18},
			{count: 2, dataWords: 19},
		}},
		"H": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 4, dataWords: 14},
			{count: 2, dataWords: 15},
		}},
	},
	9: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{{count: 2, dataWords: 116}}},
		"M": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{
			{count: 3, dataWords: 36},
			{count: 2, dataWords: 37},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 20, blocks: []ecBlock{
			{count: 4, dataWords: 16},
			{count: 4, dataWords: 17},
		}},
		"H": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 4, dataWords: 12},
			{count: 4, dataWords: 13},
		}},
	},
	10: {
		"L": ecBlocks{ecWordsPerBlock: 18, blocks: []ecBlock{
			{count: 2, dataWords: 68},
			{count: 2, dataWords: 69},
		}},
		"M": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 4, dataWords: 43},
			{count: 1, dataWords: 44},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 6, dataWords: 19},
			{count: 2, dataWords: 20},
		}},
		"H": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 6, dataWords: 15},
			{count: 2, dataWords: 16},
		}},
	},
	11: {
		"L": ecBlocks{ecWordsPerBlock: 20, blocks: []ecBlock{{count: 4, dataWords: 81}}},
		"M": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 1, dataWords: 50},
			{count: 4, dataWords: 51},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 4, dataWords: 22},
			{count: 4, dataWords: 23},
		}},
		"H": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 3, dataWords: 12},
			{count: 8, dataWords: 13},
		}},
	},
	12: {
		"L": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 2, dataWords: 92},
			{count: 2, dataWords: 93},
		}},
		"M": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{
			{count: 6, dataWords: 36},
			{count: 2, dataWords: 37},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 4, dataWords: 20},
			{count: 6, dataWords: 21},
		}},
		"H": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 7, dataWords: 14},
			{count: 4, dataWords: 15},
		}},
	},
	13: {
		"L": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{{count: 4, dataWords: 107}}},
		"M": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{
			{count: 8, dataWords: 37},
			{count: 1, dataWords: 38},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 8, dataWords: 20},
			{count: 4, dataWords: 21},
		}},
		"H": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{
			{count: 12, dataWords: 11},
			{count: 4, dataWords: 12},
		}},
	},
	14: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 3, dataWords: 115},
			{count: 1, dataWords: 116},
		}},
		"M": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 4, dataWords: 40},
			{count: 5, dataWords: 41},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 20, blocks: []ecBlock{
			{count: 11, dataWords: 16},
			{count: 5, dataWords: 17},
		}},
		"H": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 11, dataWords: 12},
			{count: 5, dataWords: 13},
		}},
	},
	15: {
		"L": ecBlocks{ecWordsPerBlock: 22, blocks: []ecBlock{
			{count: 5, dataWords: 87},
			{count: 1, dataWords: 88},
		}},
		"M": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 5, dataWords: 41},
			{count: 5, dataWords: 42},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 5, dataWords: 24},
			{count: 7, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 11, dataWords: 12},
			{count: 7, dataWords: 13},
		}},
	},
	16: {
		"L": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 5, dataWords: 98},
			{count: 1, dataWords: 99},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 7, dataWords: 45},
			{count: 3, dataWords: 46},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{
			{count: 15, dataWords: 19},
			{count: 2, dataWords: 20},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 3, dataWords: 15},
			{count: 13, dataWords: 16},
		}},
	},
	17: {
		"L": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 1, dataWords: 107},
			{count: 5, dataWords: 108},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 10, dataWords: 46},
			{count: 1, dataWords: 47},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 1, dataWords: 22},
			{count: 15, dataWords: 23},
		}},
		"H": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 2, dataWords: 14},
			{count: 17, dataWords: 15},
		}},
	},
	18: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 5, dataWords: 120},
			{count: 1, dataWords: 121},
		}},
		"M": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 9, dataWords: 43},
			{count: 4, dataWords: 44},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 17, dataWords: 22},
			{count: 1, dataWords: 23},
		}},
		"H": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 2, dataWords: 14},
			{count: 19, dataWords: 15},
		}},
	},
	19: {
		"L": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 3, dataWords: 113},
			{count: 4, dataWords: 114},
		}},
		"M": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 3, dataWords: 44},
			{count: 11, dataWords: 45},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 17, dataWords: 21},
			{count: 4, dataWords: 22},
		}},
		"H": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 9, dataWords: 13},
			{count: 16, dataWords: 14},
		}},
	},
	20: {
		"L": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 3, dataWords: 107},
			{count: 5, dataWords: 108},
		}},
		"M": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 3, dataWords: 41},
			{count: 13, dataWords: 42},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 15, dataWords: 24},
			{count: 5, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 15, dataWords: 15},
			{count: 10, dataWords: 16},
		}},
	},
	21: {
		"L": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 4, dataWords: 116},
			{count: 4, dataWords: 117},
		}},
		"M": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{{count: 17, dataWords: 42}}},
		"Q": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 17, dataWords: 22},
			{count: 6, dataWords: 23},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 19, dataWords: 16},
			{count: 6, dataWords: 17},
		}},
	},
	22: {
		"L": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 2, dataWords: 111},
			{count: 7, dataWords: 112},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{{count: 17, dataWords: 46}}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 7, dataWords: 24},
			{count: 16, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 24, blocks: []ecBlock{{count: 34, dataWords: 13}}},
	},
	23: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 4, dataWords: 121},
			{count: 5, dataWords: 122},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 4, dataWords: 47},
			{count: 14, dataWords: 48},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 11, dataWords: 24},
			{count: 14, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 16, dataWords: 15},
			{count: 14, dataWords: 16},
		}},
	},
	24: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 6, dataWords: 117},
			{count: 4, dataWords: 118},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 6, dataWords: 45},
			{count: 14, dataWords: 46},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 11, dataWords: 24},
			{count: 16, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 30, dataWords: 16},
			{count: 2, dataWords: 17},
		}},
	},
	25: {
		"L": ecBlocks{ecWordsPerBlock: 26, blocks: []ecBlock{
			{count: 8, dataWords: 106},
			{count: 4, dataWords: 107},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 8, dataWords: 47},
			{count: 13, dataWords: 48},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 7, dataWords: 24},
			{count: 22, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 22, dataWords: 15},
			{count: 13, dataWords: 16},
		}},
	},
	26: {
		"L": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 10, dataWords: 114},
			{count: 2, dataWords: 115},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 19, dataWords: 46},
			{count: 4, dataWords: 47},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 28, dataWords: 22},
			{count: 6, dataWords: 23},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 33, dataWords: 16},
			{count: 4, dataWords: 17},
		}},
	},
	27: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 8, dataWords: 122},
			{count: 4, dataWords: 123},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 22, dataWords: 45},
			{count: 3, dataWords: 46},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 8, dataWords: 23},
			{count: 26, dataWords: 24},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 12, dataWords: 15},
			{count: 28, dataWords: 16},
		}},
	},
	28: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 3, dataWords: 117},
			{count: 10, dataWords: 118},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 3, dataWords: 45},
			{count: 23, dataWords: 46},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 4, dataWords: 24},
			{count: 31, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 11, dataWords: 15},
			{count: 31, dataWords: 16},
		}},
	},
	29: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 7, dataWords: 116},
			{count: 7, dataWords: 117},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 21, dataWords: 45},
			{count: 7, dataWords: 46},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 1, dataWords: 23},
			{count: 37, dataWords: 24},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 19, dataWords: 15},
			{count: 26, dataWords: 16},
		}},
	},
	30: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 5, dataWords: 115},
			{count: 10, dataWords: 116},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 19, dataWords: 47},
			{count: 10, dataWords: 48},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 15, dataWords: 24},
			{count: 25, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 23, dataWords: 15},
			{count: 25, dataWords: 16},
		}},
	},
	31: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 13, dataWords: 115},
			{count: 3, dataWords: 116},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 2, dataWords: 46},
			{count: 29, dataWords: 47},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 42, dataWords: 24},
			{count: 1, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 23, dataWords: 15},
			{count: 28, dataWords: 16},
		}},
	},
	32: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{{count: 17, dataWords: 115}}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 10, dataWords: 46},
			{count: 23, dataWords: 47},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 10, dataWords: 24},
			{count: 35, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 19, dataWords: 15},
			{count: 35, dataWords: 16},
		}},
	},
	33: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 17, dataWords: 115},
			{count: 1, dataWords: 116},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 14, dataWords: 46},
			{count: 21, dataWords: 47},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 29, dataWords: 24},
			{count: 19, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 11, dataWords: 15},
			{count: 46, dataWords: 16},
		}},
	},
	34: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 13, dataWords: 115},
			{count: 6, dataWords: 116},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 14, dataWords: 46},
			{count: 23, dataWords: 47},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 44, dataWords: 24},
			{count: 7, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 59, dataWords: 16},
			{count: 1, dataWords: 17},
		}},
	},
	35: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 12, dataWords: 121},
			{count: 7, dataWords: 122},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 12, dataWords: 47},
			{count: 26, dataWords: 48},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 39, dataWords: 24},
			{count: 14, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 22, dataWords: 15},
			{count: 41, dataWords: 16},
		}},
	},
	36: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 6, dataWords: 121},
			{count: 14, dataWords: 122},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 6, dataWords: 47},
			{count: 34, dataWords: 48},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 46, dataWords: 24},
			{count: 10, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 2, dataWords: 15},
			{count: 64, dataWords: 16},
		}},
	},
	37: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 17, dataWords: 122},
			{count: 4, dataWords: 123},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 29, dataWords: 46},
			{count: 14, dataWords: 47},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 49, dataWords: 24},
			{count: 10, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 24, dataWords: 15},
			{count: 46, dataWords: 16},
		}},
	},
	38: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 4, dataWords: 122},
			{count: 18, dataWords: 123},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 13, dataWords: 46},
			{count: 32, dataWords: 47},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 48, dataWords: 24},
			{count: 14, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 42, dataWords: 15},
			{count: 32, dataWords: 16},
		}},
	},
	39: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 20, dataWords: 117},
			{count: 4, dataWords: 118},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 40, dataWords: 47},
			{count: 7, dataWords: 48},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 43, dataWords: 24},
			{count: 22, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 10, dataWords: 15},
			{count: 67, dataWords: 16},
		}},
	},
	40: {
		"L": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 19, dataWords: 118},
			{count: 6, dataWords: 119},
		}},
		"M": ecBlocks{ecWordsPerBlock: 28, blocks: []ecBlock{
			{count: 18, dataWords: 47},
			{count: 31, dataWords: 48},
		}},
		"Q": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 34, dataWords: 24},
			{count: 34, dataWords: 25},
		}},
		"H": ecBlocks{ecWordsPerBlock: 30, blocks: []ecBlock{
			{count: 20, dataWords: 15},
			{count: 61, dataWords: 16},
		}},
	},
}
