package src

import "encoding/json"

// Compound for get net cachekey
type Compound struct {
	Cid                    int     `json:"cid,omitempty"`
	Mw                     float64 `json:"mw,omitempty"`
	Polararea              float64 `json:"polararea,omitempty"`
	Complexity             float64 `json:"complexity,omitempty"`
	Xlogp                  float64 `json:"xlogp,omitempty"`
	Heavycnt               int     `json:"heavycnt,omitempty"`
	Hbonddonor             int     `json:"hbonddonor,omitempty"`
	Hbondacc               int     `json:"hbondacc,omitempty"`
	Rotbonds               int     `json:"rotbonds,omitempty"`
	Annothitcnt            int     `json:"annothitcnt,omitempty"`
	Charge                 int     `json:"charge,omitempty"`
	Covalentunitcnt        int     `json:"covalentunitcnt,omitempty"`
	Isotopeatomcnt         int     `json:"isotopeatomcnt,omitempty"`
	Totalatomstereocnt     int     `json:"totalatomstereocnt,omitempty"`
	Definedatomstereocnt   int     `json:"definedatomstereocnt,omitempty"`
	Undefinedatomstereocnt int     `json:"undefinedatomstereocnt,omitempty"`
	Totalbondstereocnt     int     `json:"totalbondstereocnt,omitempty"`
	Definedbondstereocnt   int     `json:"definedbondstereocnt,omitempty"`
	Undefinedbondstereocnt int     `json:"undefinedbondstereocnt,omitempty"`
	Pclidcnt               int     `json:"pclidcnt,omitempty"`
	Gpidcnt                int     `json:"gpidcnt,omitempty"`
	Gpfamilycnt            int     `json:"gpfamilycnt,omitempty"`
	Aids                   string  `json:"aids,omitempty"`
	Cmpdname               string  `json:"cmpdname,omitempty"`
	Cmpdsynonym            string  `json:"cmpdsynonym,omitempty"`
	Inchi                  string  `json:"inchi,omitempty"`
	Inchikey               string  `json:"inchikey,omitempty"`
	Smiles                 string  `json:"smiles,omitempty"`
	Iupacname              string  `json:"iupacname,omitempty"`
	Mf                     string  `json:"mf,omitempty"`
	Sidsrcname             string  `json:"sidsrcname,omitempty"`
	Annotation             string  `json:"annotation,omitempty"`
	Cidcdate               string  `json:"cidcdate,omitempty"`
	Depcatg                string  `json:"depcatg,omitempty"`
	Meshheadings           string  `json:"meshheadings,omitempty"`
	Annothits              string  `json:"annothits,omitempty"`
	Exactmass              string  `json:"exactmass,omitempty"`
	Monoisotopicmass       string  `json:"monoisotopicmass,omitempty"`
}

type SearchType struct {
	Cid      string `json:"Cid"`
	Smiles   string `json:"Smiles"`
	Name     string `json:"Name"`
	Inchi    string `json:"Inchi"`
	Inchikey string `json:"Inchikey"`
}

type PubchemCache struct {
	Response struct {
		Status            int       `json:"status"`
		Message           []*string `json:"message"`
		Hitcount          int       `json:"hitcount"`
		Percentcompletion *float64  `json:"percentcompletion"`
		Cachekey          string    `json:"cachekey"`
	} `json:"response"`
}

// Compounds for pug api
type Compounds struct {
	PCCompounds []struct {
		Id struct {
			Id struct {
				Cid int `json:"cid"`
			} `json:"id"`
		} `json:"id"`
		Atoms struct {
			Aid     []int `json:"aid"`
			Element []int `json:"element"`
		} `json:"atoms"`
		Bonds struct {
			Aid1  []int `json:"aid1"`
			Aid2  []int `json:"aid2"`
			Order []int `json:"order"`
		} `json:"bonds"`
		Coords []struct {
			Type       []int `json:"type"`
			Aid        []int `json:"aid"`
			Conformers []struct {
				X     []float64 `json:"x"`
				Y     []float64 `json:"y"`
				Style struct {
					Annotation []int `json:"annotation"`
					Aid1       []int `json:"aid1"`
					Aid2       []int `json:"aid2"`
				} `json:"style"`
			} `json:"conformers"`
		} `json:"coords"`
		Charge int `json:"charge"`
		Props  []struct {
			Urn struct {
				Label          string `json:"label"`
				Name           string `json:"name,omitempty"`
				Datatype       int    `json:"datatype"`
				Release        string `json:"release"`
				Implementation string `json:"implementation,omitempty"`
				Version        string `json:"version,omitempty"`
				Software       string `json:"software,omitempty"`
				Source         string `json:"source,omitempty"`
				Parameters     string `json:"parameters,omitempty"`
			} `json:"urn"`
			Value struct {
				Ival   int     `json:"ival,omitempty"`
				Fval   float64 `json:"fval,omitempty"`
				Binary string  `json:"binary,omitempty"`
				Sval   string  `json:"sval,omitempty"`
			} `json:"value"`
		} `json:"props"`
		Count struct {
			HeavyAtom       int `json:"heavy_atom"`
			AtomChiral      int `json:"atom_chiral"`
			AtomChiralDef   int `json:"atom_chiral_def"`
			AtomChiralUndef int `json:"atom_chiral_undef"`
			BondChiral      int `json:"bond_chiral"`
			BondChiralDef   int `json:"bond_chiral_def"`
			BondChiralUndef int `json:"bond_chiral_undef"`
			IsotopeAtom     int `json:"isotope_atom"`
			CovalentUnit    int `json:"covalent_unit"`
			Tautomers       int `json:"tautomers"`
		} `json:"count"`
	} `json:"PC_Compounds"`
}

type pubChemError struct {
	ErrorCode int
	ErrorMsg  string
}

type CompoundProperty struct {
	MolecularFormula         string // Molecular formula.
	MolecularWeight          string // The molecular weight is the sum of all atomic weights of the constituent atoms in a compound, measured in g/mol.In the absence of explicit isotope labelling, averaged natural abundance is assumed.If an atom bears an explicit isotope label, 100% isotopic purity is assumed at this location.
	CanonicalSMILES          string // Canonical SMILES (Simplified Molecular Input Line Entry System) string.It is a unique SMILES string of a compound, generated by a “canonicalization” algorithm.
	IsomericSMILES           string // Isomeric SMILES string.It is a SMILES string with stereochemical and isotopic specifications.
	InChI                    string // Standard IUPAC International Chemical Identifier (InChI).It does not allow for user selectable options in dealing with the stereochemistry and tautomer layers of the InChI string.
	InChIKey                 string // Hashed version of the full standard InChI, consisting of 27 characters.
	IUPACName                string // Chemical name systematically determined according to the IUPAC nomenclatures.
	Title                    string // The title used for the compound summary page.
	XLogP                    string // Computationally generated octanol-water partition coefficient or distribution coefficient.XLogP is used as a measure of hydrophilicity or hydrophobicity of a molecule.
	ExactMass                string // The mass of the most likely isotopic composition for a single molecule, corresponding to the most intense ion/molecule peak in a mass spectrum.
	MonoisotopicMass         string // The mass of a molecule, calculated using the mass of the most abundant isotope of each element.
	TPSA                     string // Topological polar surface area, computed by the algorithm described in the paper by Ertl et al.
	Complexity               string // The molecular complexity rating of a compound, computed using the Bertz/Hendrickson/Ihlenfeldt formula.
	Charge                   string // The total (or net) charge of a molecule.
	HBondDonorCount          string // Number of hydrogen-bond donors in the structure.
	HBondAcceptorCount       string // Number of hydrogen-bond acceptors in the structure.
	RotatableBondCount       string // Number of rotatable bonds.
	HeavyAtomCount           string // Number of non-hydrogen atoms.
	IsotopeAtomCount         string // Number of atoms with enriched isotope(s)
	AtomStereoCount          string // Total number of atoms with tetrahedral (sp3) stereo [e.g., (R)- or (S)-configuration]
	DefinedAtomStereoCount   string // Number of atoms with defined tetrahedral (sp3) stereo.
	UndefinedAtomStereoCount string // Number of atoms with undefined tetrahedral (sp3) stereo.
	BondStereoCount          string // Total number of bonds with planar (sp2) stereo [e.g., (E)- or (Z)-configuration].
	DefinedBondStereoCount   string // Number of atoms with defined planar (sp2) stereo.
	UndefinedBondStereoCount string // Number of atoms with undefined planar (sp2) stereo.
	CovalentUnitCount        string // Number of covalently bound units.
	PatentCount              string // Number of patent documents linked to this compound.
	PatentFamilyCount        string // Number of unique patent families linked to this compound (e.g.patent documents grouped by family).
	LiteratureCount          string //  Number of articles linked to this compound (by PubChem's consolidated literature analysis).
	Volume3D                 string //  Analytic volume of the first diverse conformer (default conformer) for a compound.
	XStericQuadrupole3D      string // The x component of the quadrupole moment (Qx) of the first diverse conformer (default conformer) for a compound.
	YStericQuadrupole3D      string // The y component of the quadrupole moment (Qy) of the first diverse conformer (default conformer) for a compound.
	ZStericQuadrupole3D      string // The z component of the quadrupole moment (Qz) of the first diverse conformer (default conformer) for a compound.
	FeatureCount3D           string // Total number of 3D features (the sum of FeatureAcceptorCount3D, FeatureDonorCount3D, FeatureAnionCount3D, FeatureCationCount3D, FeatureRingCount3D and FeatureHydrophobeCount3D)
	FeatureAcceptorCount3D   string //  Number of hydrogen-bond acceptors of a conformer.
	FeatureDonorCount3D      string //  Number of hydrogen-bond donors of a conformer.
	FeatureAnionCount3D      string //  Number of anionic centers (at pH 7) of a conformer.
	FeatureCationCount3D     string // Number of cationic centers (at pH 7) of a conformer.
	FeatureRingCount3D       string // Number of rings of a conformer.
	FeatureHydrophobeCount3D string // Number of hydrophobes of a conformer.
	ConformerModelRMSD3D     string // Conformer sampling RMSD in Å.
	EffectiveRotorCount3D    string // Total number of 3D features (the sum of FeatureAcceptorCount3D, FeatureDonorCount3D, FeatureAnionCount3D, FeatureCationCount3D, FeatureRingCount3D and FeatureHydrophobeCount3D)
	ConformerCount3D         string // The number of conformers in the conformer model for a compound.
	Fingerprint2D            string //  Base64-encoded PubChem Substructure Fingerprint of a molecule.
}

type NetCacheKeyPayload struct {
	Select     string `json:"select"`
	Collection string `json:"collection"`
	Where      struct {
		Ands []struct {
			Input struct {
				Type   string `json:"type"`
				Idtype string `json:"idtype"`
				Key    string `json:"key"`
			} `json:"input"`
		} `json:"ands"`
	} `json:"where"`
	Order   []string `json:"order"`
	Start   int      `json:"start"`
	Limit   int      `json:"limit"`
	Width   int      `json:"width"`
	Listids int      `json:"listids"`
}

type SDQSet struct {
	Status struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	} `json:"status"`
	InputCount int        `json:"inputCount"`
	TotalCount int        `json:"totalCount"`
	Collection string     `json:"collection"`
	Type       string     `json:"type"`
	Rows       []Compound `json:"rows"`
}

// SDQOutputSet for search from cachekey
type SDQOutputSet struct {
	SDQOutputSet []SDQSet `json:"SDQOutputSet"`
}

type Parameter struct {
	Name   string `json:"name"`
	String string `json:"string,omitempty"`
	Bool   bool   `json:"bool,omitempty"`
	Num    int    `json:"num,omitempty"`
}

type QueryBlob struct {
	Query struct {
		Type      string      `json:"type"`
		Parameter []Parameter `json:"parameter"`
	} `json:"query"`
}

func (v QueryBlob) toString() string {
	queryJs, _ := json.Marshal(v)
	s := string(queryJs)
	return s
}

func (c Compounds) Get() {

}
