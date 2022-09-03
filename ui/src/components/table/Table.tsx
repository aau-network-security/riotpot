import "./Table.scss";

export const Cell = ({ content }: { content: any }) => {
  return <td className="tableCell">{content}</td>;
};

const Header = ({ content }: { content: any }) => {
  return <th className="tableCell">{content}</th>;
};

const Headers = ({ headers }: { headers: any }) => {
  return (
    <thead>
      <tr className="tableHeaders">
        {headers.map((cell: any) => {
          return <Header content={cell} />;
        })}
      </tr>
    </thead>
  );
};

export const Row = ({ cells }: { cells: any }) => {
  return (
    <tr className="tableRow">
      {cells.map((content: any, ind: Number) => {
        return <Cell content={content} />;
      })}
    </tr>
  );
};

const Body = ({ rows, children }: { rows: any[]; children: any[] }) => {
  return (
    <tbody>
      {rows.map((cells: any[], ind: Number) => {
        return <Row cells={cells} />;
      })}
      {children}
    </tbody>
  );
};

export const Table = ({ data, rows }: { data: any; rows?: any }) => {
  return (
    <div className="table-container">
      <table>
        <Headers headers={data.headers} />
        <Body rows={data.rows}>{rows}</Body>
      </table>
    </div>
  );
};

export default Table;
