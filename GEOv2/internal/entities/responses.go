package entities

type SearchResponse struct {
    Addresses []Address
}

type Address struct {
    Source       string  
    Result       string  
    PostalCode   string  
    Country      string  
    Region       string 
    CityArea     string 
    CityDistrict string 
    Street       string 
    House        string  
    GeoLat       float64
    GeoLon       float64  
}

type GeocodeResponse struct {
    Suggestions []Suggestion
}

type Suggestion struct {
    Value             string
    UnrestrictedValue string
    Data              SuggestionData
}

type SuggestionData struct {
    PostalCode   string
    Country      string
    Region       string
    CityArea     string
    CityDistrict string
    Street       string
    House        string
    GeoLat       float64
    GeoLon       float64
}