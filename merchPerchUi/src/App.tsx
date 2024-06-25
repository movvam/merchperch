// import { useState } from 'react'
// import { Button } from "@/components/ui/button"
import ArtistJsonDataDisplay from "@/components/ui/artistJsonTable"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"


function App() {
  // const [count, setCount] = useState(0)

  return (
    <>
      {/* <Button>Click me</Button> */}
      <div className="app">
      <header className="App-header">
        <h1>
          MerchPerch 
        </h1>
        <h5>
          Find merch for all your favorite Spotify artists
        </h5>
      </header>
      <ArtistJsonDataDisplay/>

    </div>
      
    </>
  )
}

export default App
