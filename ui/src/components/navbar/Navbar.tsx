import Nav from "react-bootstrap/Nav";
import "./Navbar.scss";
import { useLocation } from "react-router-dom";
import { Pages } from "../../constants/globals";

const Navbar = () => {
  const location = useLocation();
  const { pathname } = location;
  const curr = pathname.split("/")[1];
  const def = "/instance";

  return (
    <Nav
      activeKey={curr ? "/" + curr : def}
      defaultActiveKey={def}
      className="flex-column"
    >
      <h5>Navigation</h5>
      {Pages.map((obj) => {
        // Convert the page title into lower case
        const link = "/" + obj.page.toLowerCase();

        // Return a link object
        return (
          <Nav.Link key={obj.page} href={link}>
            <obj.icon />
            {obj.page}
          </Nav.Link>
        );
      })}
    </Nav>
  );
};

export default Navbar;
