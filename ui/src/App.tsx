import { Routes, Route } from "react-router-dom";

/* CSS */
import { Container, Col, Row } from "react-bootstrap";
import "./App.scss";

/* Components */
import Navbar from "./components/navbar/Navbar";
import Services from "./routes/services/Services";
import Instances from "./routes/instances/Instances";
import Profiles from "./routes/profiles/Profiles";

/* Pages */

function App() {
  return (
    <div className="App">
      <Container>
        <Row>
          {/* Navigation */}
          <Col>
            <Navbar />
          </Col>

          {/* Content */}
          <Col xs={7}>
            <Routes>
              <Route path="instances" element={<Instances />}>
                <Route path=":id"></Route>
              </Route>
              <Route path="services" element={<Services />}>
                <Route path=":id"></Route>
              </Route>
              <Route path="profiles" element={<Profiles />}>
                <Route path=":id"></Route>
              </Route>
              <Route path="settings"></Route>
              <Route
                path="*"
                element={
                  <main style={{ padding: "1rem" }}>
                    <p>Page not found</p>
                  </main>
                }
              ></Route>
            </Routes>
          </Col>

          {/* Right space */}
          <Col></Col>
        </Row>
      </Container>
    </div>
  );
}

export default App;
