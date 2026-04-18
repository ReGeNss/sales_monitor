package regexp

const GramsRegex = `(\d+(?:[.,]\d+)?)\s*(грам|гр|г)\.?`

const KilogramRegex = `(\d+(?:[.,]\d+)?)\s*(кг|кілограм|кіло|кг|кіло)\.?`

const VolumeLiterRegex = `(\d+(?:[.,]\d+)?)\s*(л|літр[а-я]*)\.?`

const VolumeMilliliterRegex = `(\d+(?:[.,]\d+)?)\s*(мл|млітри|млітра)\.?`

const WithoutDecimalRegex = `[\d.,]`
