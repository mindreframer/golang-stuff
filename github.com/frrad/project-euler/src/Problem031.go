package main

import "fmt"

func main() {
	count := 0

	for twobucks := 0; twobucks <= 1; twobucks++ {

		for abuck := 0; abuck <= 2; abuck++ {

			if twobucks*200+abuck*100 <= 200 {
				for fiddy := 0; fiddy <= 4; fiddy++ {

					if twobucks*200+abuck*100+fiddy*50 <= 200 {
						for twenty := 0; twenty <= 10; twenty++ {

							if twobucks*200+abuck*100+fiddy*50+twenty*20 <= 200 {
								for ten := 0; ten <= 20; ten++ {

									if twobucks*200+abuck*100+fiddy*50+twenty*20+ten*10 <= 200 {
										for five := 0; five <= 40; five++ {

											if twobucks*200+abuck*100+fiddy*50+twenty*20+ten*10+five*5 <= 200 {
												for two := 0; two <= 100; two++ {

													if twobucks*200+abuck*100+fiddy*50+twenty*20+ten*10+five*5+two*2 <= 200 {
														for one := 0; one <= 200; one++ {

															if twobucks*200+abuck*100+fiddy*50+twenty*20+ten*10+five*5+two*2+one == 200 {
																count++
															}

														}
													}
												}

											}

										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	fmt.Println(count)
}
