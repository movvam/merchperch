import ProductJsonDataDisplay from "@/components/productsJsonDataDisplay";

function ProductsPage() {
  return (
    <div className="app">
      <header className="App-header">
        <h1>
          MerchPerch 
        </h1>
        <h5>
          Find merch for all your favorite Spotify artists
        </h5>
      </header>
      <ProductJsonDataDisplay/>
    </div>
    )
}

export default ProductsPage

