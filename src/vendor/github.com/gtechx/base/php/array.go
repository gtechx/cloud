package php

/* ============================================================================================ */
func InArrayFloat64(arr []float64, needle float64) bool {
    for _, item := range arr {
        if item == needle {
            return true
        }
    }
    return false
}
/* ============================================================================================ */
func InArrayInt64(arr []int64, needle int64) bool {
    for _, item := range arr {
        if item == needle {
            return true
        }
    }
    return false
}
/* ============================================================================================ */
func InArrayInt(arr []int, needle int) bool {
    for _, item := range arr {
        if item == needle {
            return true
        }
    }
    return false
}
/* ============================================================================================ */
func InArrayString(arr []string, needle string) bool {
    for _, item := range arr {
        if item == needle {
            return true
        }
    }
    return false
}
/* ============================================================================================ */
func ArraySearchString(arr []string, needle string) int {
    for i, item := range arr {
        if item == needle {
            return i
        }
    }
    return -1
}
/* ============================================================================================ */
func ArraySearchInt(arr []int, needle int) int {
    for i, item := range arr {
        if item == needle {
            return i
        }
    }
    return -1
}
/* ============================================================================================ */
func ArraySearchInt64(arr []int64, needle int64) int {
    for i, item := range arr {
        if item == needle {
            return i
        }
    }
    return -1
}