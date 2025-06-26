import React from 'react';
import type { Bill } from './types';

interface BillListProps {
  bills: Bill[];
}

const BillList: React.FC<BillListProps> = ({ bills }) => {
  return (
    <div className="bills-list">
      {bills.map((bill) => (
        <div key={bill.legisinfo_id} className="bill-item">
          <h2>{bill.number}: {bill.name.en}</h2>
          <p>Session: {bill.session}</p>
          <p>Introduced: {new Date(bill.introduced).toLocaleDateString()}</p>
          <p><a href={bill.url} target="_blank" rel="noopener noreferrer">View Bill (Legisinfo)</a></p>
        </div>
      ))}
    </div>
  );
};

export default BillList;
