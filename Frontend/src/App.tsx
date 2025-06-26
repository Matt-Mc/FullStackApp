import React, { useState, useEffect } from 'react';
import './App.css';

// --- Define TypeScript Interfaces (ONLY for the Bill structure) ---
interface BillName {
  EN: string;
  FR: string;
}

interface Bill {
  Session: string;
  Legisinfo_id: number;
  Introduced: string;
  Name: BillName;
  Number: string;
  Url: string;
}
// --- End of Interfaces ---


function App() {
  // State to store the bills data, explicitly typed as an array of Bill
  const [bills, setBills] = useState<Bill[]>([]);

  // State to manage loading status
  const [loading, setLoading] = useState<boolean>(true);

  // State to manage potential errors
  const [error, setError] = useState<Error | null>(null);

  // useEffect hook to fetch data when the component mounts
  useEffect(() => {
    const fetchBills = async () => {
      try {
        const response = await fetch('http://localhost:1500/api/bills'); // Your Go backend endpoint
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        // --- CRITICAL CHANGE HERE: Assign directly to data ---
        // Expecting the response to be a direct array of Bill objects
        const data: Bill[] = await response.json(); // Type assertion is now correct

        setBills(data); // Set the directly received array to your state
        // --- END CRITICAL CHANGE ---

      } catch (err) {
        if (err instanceof Error) {
          setError(err);
        } else {
          setError(new Error('An unknown error occurred'));
        }
      } finally {
        setLoading(false);
      }
    };

    fetchBills();
  }, []);

  if (loading) {
    return <div className="App">Loading bills...</div>;
  }

  if (error) {
    return (
      <div className="App">
        <h1>Error: {error.message}</h1>
        <p>Could not load bills. Please ensure your Go backend is running on http://localhost:8080</p>
      </div>
    );
  }

  return (
    <div className="App">
      <h1>Canadian Bills</h1>
      <div className="bills-list">
        {bills.map((bill: Bill) => (
          <div key={bill.Legisinfo_id} className="bill-item">
            <h2>{bill.Number}: {bill.Name.EN}</h2>
            <p>Session: {bill.Session}</p>
            <p>Introduced: {new Date(bill.Introduced).toLocaleDateString()}</p>
            <p><a href={bill.Url} target="_blank" rel="noopener noreferrer">View Bill (Legisinfo)</a></p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;