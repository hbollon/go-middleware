// BOLLON hugo / RODRIGUEZ Samuel
package main

type StateSectionCritique int

const (
	StateSectionCritique_Released StateSectionCritique = iota
	StateSectionCritique_Requested
	StateSectionCritique_Critical_Section
)
