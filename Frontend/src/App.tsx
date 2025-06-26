import React, { useState, useEffect } from 'react';
import './App.css';
import { fetchBills, fetchMPs } from './api';
import type {Bill, MP } from './types';
import BillList from './BillList';
import MPList from './MPList';
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer } from 'recharts';

const PARTY_COLORS: Record<string, string> = {
  "Liberal": "#D71A4F",
  "Conservative": "#1A4782",
  "NDP": "#FF6600",
  "Bloc Québécois": "#008080",
  "Green Party": "#008000",
  "Independent": "#FFFFFF",
  // Add more parties and colors as needed
};

function App() {
  const [bills, setBills] = useState<Bill[]>([]);
  const [mps, setMPs] = useState<MP[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);
  const [partyTally, setPartyTally] = useState<Record<string, number>>({});

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [billsData, mpsData] = await Promise.all([
          fetchBills(),
          fetchMPs(),
        ]);
        setBills(billsData);
        setMPs(mpsData);

        const tally: Record<string, number> = {};
        mpsData.forEach(mp => {
          const partyName = mp.current_party.short_name.en;
          tally[partyName] = (tally[partyName] || 0) + 1;
        });
        setPartyTally(tally);

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

    fetchData();
  }, []);

  if (loading) {
    return <div className="App">Loading...</div>;
  }

  if (error) {
    return (
      <div className="App">
        <h1>Error: {error.message}</h1>
        <p>Could not load data. Please ensure your Go backend is running on http://localhost:1500</p>
      </div>
    );
  }

  const pieChartData = Object.entries(partyTally).map(([name, value]) => ({
    name,
    value,
  }));

  return (
    <div className="App">
      <h1>Canadian Parliament Dashboard</h1>
      <div className="party-tally-section">
        <h2>Party Distribution</h2>
        <div className="chart-container">
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={pieChartData}
                cx="50%"
                cy="50%"
                outerRadius={100}
                fill="#8884d8"
                dataKey="value"
                label
              >
                {pieChartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={PARTY_COLORS[entry.name] || "#CCCCCC"} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
        <div className="tally-list">
          <ul>
            {Object.entries(partyTally).map(([party, count]) => (
              <li key={party} style={{ color: PARTY_COLORS[party] || "#000000" }}>
                {party}: {count} MPs
              </li>
            ))}
          </ul>
        </div>
      </div>
      <div className="data-container">
        <div className="bills-container">
          <h2>Bills</h2>
          <BillList bills={bills} />
        </div>
        <div className="mps-container">
          <h2>Members of Parliament</h2>
          <MPList mps={mps} />
        </div>
      </div>
    </div>
  );
}

export default App;