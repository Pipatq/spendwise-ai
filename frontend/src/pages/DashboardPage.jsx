import React, { useEffect, useState } from 'react';
import { getSpendingSummary, generateSummary } from '../services/api';

const DashboardPage = () => {
  const [spendingData, setSpendingData] = useState([]);
  const [summary, setSummary] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await getSpendingSummary();
        setSpendingData(response.data);
      } catch (err) {
        setError('Failed to fetch spending data');
      }
    };
    fetchData();
  }, []);

  const handleGenerateSummary = async () => {
    try {
      const response = await generateSummary();
      setSummary(response.data.summary);
    } catch (err) {
      setError('Failed to generate summary');
    }
  };

  return (
    <div className="container p-4 mx-auto">
      <h2 className="text-2xl font-bold">Dashboard</h2>
      {error && <p className="text-red-500">{error}</p>}
      <div className="mt-4">
        <h3 className="text-xl font-semibold">Spending Summary</h3>
        <ul className="mt-2">
          {spendingData.map((item, index) => (
            <li key={index} className="flex justify-between py-1">
              <span>{item.category}</span>
              <span>${item.amount.toFixed(2)}</span>
            </li>
          ))}
        </ul>
      </div>
      <div className="mt-6">
        <button
          onClick={handleGenerateSummary}
          className="px-4 py-2 text-white bg-blue-600 rounded-md"
        >
          Generate AI Summary
        </button>
        {summary && (
          <div className="p-4 mt-4 bg-gray-100 rounded-md">
            <h4 className="font-semibold">AI Summary:</h4>
            <p>{summary}</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default DashboardPage;
