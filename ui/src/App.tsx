import { Routes, Route } from "react-router-dom";

/* CSS */
import { Container, Col, Row } from "react-bootstrap";
import "./App.scss";

/* Components */
import Navbar from "./components/navbar/Navbar";
import Services from "./routes/services/Services";
import Profiles from "./routes/profiles/Profiles";
import Profile from "./routes/profiles/Profile";
import Instance from "./routes/instances/Instance";
import { SimpleBreadcrumb } from "./components/utils/Common";
import { Settings } from "./routes/settings/Settings";

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
          <Col xs={9}>
            <SimpleBreadcrumb />
            <Routes>
              <Route path="instance" element={<Instance />}></Route>
              <Route path="services">
                <Route index element={<Services />} />
                <Route path=":id"></Route>
              </Route>
              <Route path="profiles">
                <Route index element={<Profiles />} />
                <Route path=":id" element={<Profile />}></Route>
              </Route>
              <Route path="settings" element={<Settings />}></Route>
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
