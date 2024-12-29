import React, { useEffect, useState, useRef } from 'react';
import { fetchPaymentPostResponses, PaymentPostResponse } from '../api/paymentPostResponse';

const PaymentPostResponseTable: React.FC = () => {
  const [responses, setResponses] = useState<PaymentPostResponse[]>([]);
  const [filteredResponses, setFilteredResponses] = useState<PaymentPostResponse[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [rowsPerPage] = useState<number>(10); // Showing 10 rows per page
  const [searchQuery, setSearchQuery] = useState<string>('');
  const topPaginationRef = useRef<HTMLDivElement>(null);
  const bottomPaginationRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await fetchPaymentPostResponses();
        setResponses(data);
        setFilteredResponses(data);
        setLoading(false);
      } catch (err) {
        setError('Failed to load payment post responses');
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  // Search and filter logic
  useEffect(() => {
    const filteredData = responses.filter((response) =>
      response.payer_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      response.pprno?.toString().toLowerCase().includes(searchQuery.toLowerCase()) || // Convert number to string
      response.applicant_id?.toString().toLowerCase().includes(searchQuery.toLowerCase()) // Convert number to string
    );
    setFilteredResponses(filteredData);
  }, [searchQuery, responses]);

  // Calculate pagination
  const indexOfLastRow = currentPage * rowsPerPage;
  const indexOfFirstRow = indexOfLastRow - rowsPerPage;
  const currentResponses = filteredResponses.slice(indexOfFirstRow, indexOfLastRow);
  const totalPages = Math.ceil(filteredResponses.length / rowsPerPage);

  const maxPageNumbersToShow = 5;

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const getPaginationGroup = () => {
    let start = Math.floor((currentPage - 1) / maxPageNumbersToShow) * maxPageNumbersToShow;
    return new Array(Math.min(maxPageNumbersToShow, totalPages - start))
      .fill(null)
      .map((_, index) => start + index + 1);
  };

  // Hide top pagination when bottom pagination is visible
  const handleScroll = () => {
    if (bottomPaginationRef.current && topPaginationRef.current) {
      const bottomPosition = bottomPaginationRef.current.getBoundingClientRect().top;
      const topPosition = topPaginationRef.current.getBoundingClientRect().bottom;

      if (bottomPosition <= window.innerHeight && topPosition < bottomPosition) {
        topPaginationRef.current.style.display = 'none';
      } else {
        topPaginationRef.current.style.display = 'flex';
      }
    }
  };

  useEffect(() => {
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  if (loading) {
    return <div className="text-center p-5">Loading...</div>;
  }

  if (error) {
    return <div className="text-center text-red-500 p-5">{error}</div>;
  }

  return (
    <div className="p-4">

      <div className="sticky top-0 z-20 bg-transparent">
        <div className="bg-[#EEF0F1]"> {/* Apply background to this wrapper div */}
          {/* Payment Post Responses heading */}
          <h1 className="text-[18px] font-bold mb-4 text-center text-primary sticky top-0 z-20">
            Payment Post Responses
          </h1>

          {/* Search bar */}
          <div className="mb-6 flex justify-center sticky top-[38px] z-10">
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search by PPR No, Applicant ID, or Payer Name"
              className="w-[306px] md:w-1/2 p-2 border border-gray-300 rounded-lg shadow focus:outline-none focus:ring-2 focus:ring-primary text-[13px] md:text-base bg-transparent"
            />
          </div>

          {/* Pagination at the top for mobile */}
          <div className="md:hidden mt-4 mb-4 flex justify-center space-x-2 items-center sticky top-[76px] z-10">
            {/* Double left arrow */}
            <button
              onClick={() => handlePageChange(1)}
              disabled={currentPage === 1}
              className="px-2 py-1 rounded-lg bg-gray-200 hover:bg-primary hover:text-white disabled:bg-gray-300"
            >
              &laquo;
            </button>

            {/* Single left arrow */}
            <button
              onClick={() => handlePageChange(currentPage - 1)}
              disabled={currentPage === 1}
              className="px-2 py-1 rounded-lg bg-gray-200 hover:bg-primary hover:text-white disabled:bg-gray-300"
            >
              &lt;
            </button>

            {getPaginationGroup().map((page) => (
              <button
                key={page}
                onClick={() => handlePageChange(page)}
                className={`px-2 py-1 rounded-lg ${
                  currentPage === page ? 'bg-primary text-white' : 'bg-gray-200'
                } hover:bg-primary hover:text-white`}
              >
                {page}
              </button>
            ))}

            {/* Single right arrow */}
            <button
              onClick={() => handlePageChange(currentPage + 1)}
              disabled={currentPage === totalPages}
              className="px-2 py-1 rounded-lg bg-gray-200 hover:bg-primary hover:text-white disabled:bg-gray-300"
            >
              &gt;
            </button>

            {/* Double right arrow */}
            <button
              onClick={() => handlePageChange(totalPages)}
              disabled={currentPage === totalPages}
              className="px-2 py-1 rounded-lg bg-gray-200 hover:bg-primary hover:text-white disabled:bg-gray-300"
            >
              &raquo;
            </button>
          </div>
        </div>
      </div>

      <div className="overflow-x-auto px-2 md:px-0"> {/* Added padding on mobile */}
        <table className="min-w-full border-collapse block md:table bg-white rounded-lg shadow">
          {/* Remove table headers for mobile view */}
          <thead className="hidden md:table-header-group bg-secondary text-white">
            <tr className="border border-gray-300 md:border-none block md:table-row">
              <th className="block md:table-cell text-left p-2 text-[13px]">PPR No</th>
              <th className="block md:table-cell text-left p-2 text-[13px]">Applicant ID</th>
              <th className="block md:table-cell text-left p-2 text-[13px]">Receipt</th>
              <th className="block md:table-cell text-left p-2 text-[13px]">Status</th>
              <th className="block md:table-cell text-left p-2 text-[13px]">Amount (TZS)</th>
              <th className="block md:table-cell text-left p-2 text-[13px]">Balance (TZS)</th>
              <th className="block md:table-cell text-left p-2 text-[13px]">Payer Name</th>
            </tr>
          </thead>
          <tbody className="block md:table-row-group">
            {currentResponses.map((response) => (
              <tr key={response.pprno} className="border border-gray-300 md:border-none block md:table-row">
                <td className="block md:table-cell text-left p-2 text-[13px]">
                  <span className="md:hidden font-bold text-secondary">PPR No: </span>{response.pprno}
                </td>
                <td className="block md:table-cell text-left p-2 text-[13px]">
                  <span className="md:hidden font-bold text-secondary">Applicant ID: </span>{response.applicant_id}
                </td>
                <td className="block md:table-cell text-left p-2 text-[13px]">
                  <span className="md:hidden font-bold text-secondary">Receipt: </span>{response.receipt || 'N/A'}
                </td>
                <td className="block md:table-cell text-left p-2 text-[13px]">
                  <span className="md:hidden font-bold text-secondary">Status: </span>{response.status || 'N/A'}
                </td>
                <td className="block md:table-cell text-left p-2 text-[13px]">
                  <span className="md:hidden font-bold text-secondary">Amount (TZS): </span>{response.amount_tzs ?? 'N/A'}
                </td>
                <td className="block md:table-cell text-left p-2 text-[13px]">
                  <span className="md:hidden font-bold text-secondary">Balance (TZS): </span>{response.balance_tzs ?? 'N/A'}
                </td>
                <td className="block md:table-cell text-left p-2 text-[13px]">
                  <span className="md:hidden font-bold text-secondary">Payer Name: </span>{response.payer_name || 'N/A'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>


      {/* Pagination at the bottom */}
      <div className="mt-4 mb-4 flex justify-center space-x-2 items-center" ref={bottomPaginationRef}>
        {/* Double left arrow */}
        <button
          onClick={() => handlePageChange(1)}
          disabled={currentPage === 1}
          className="px-2 py-1 rounded-lg bg-gray-200 hover:bg-primary hover:text-white disabled:bg-gray-300"
        >
          &laquo;
        </button>

        {/* Single left arrow */}
        <button
          onClick={() => handlePageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className="px-2 py-1 rounded-lg bg-gray-200 hover:bg-primary hover:text-white disabled:bg-gray-300"
        >
          &lt;
        </button>

        {getPaginationGroup().map((page) => (
          <button
            key={page}
            onClick={() => handlePageChange(page)}
            className={`px-2 py-1 rounded-lg ${
              currentPage === page ? 'bg-primary text-white' : 'bg-gray-200'
            } hover:bg-primary hover:text-white`}
          >
            {page}
          </button>
        ))}

        {/* Single right arrow */}
        <button
          onClick={() => handlePageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="px-2 py-1 rounded-lg bg-gray-200 hover:bg-primary hover:text-white disabled:bg-gray-300"
        >
          &gt;
        </button>

        {/* Double right arrow */}
        <button
          onClick={() => handlePageChange(totalPages)}
          disabled={currentPage === totalPages}
          className="px-2 py-1 rounded-lg bg-gray-200 hover:bg-primary hover:text-white disabled:bg-gray-300"
        >
          &raquo;
        </button>
      </div>
    </div>
  );
};

export default PaymentPostResponseTable;
