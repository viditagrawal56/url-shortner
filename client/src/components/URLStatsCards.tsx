interface URLStatsCardProps {
  label: string;
  count: number;
  Icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
}

const URLStatsCard = ({ label, count, Icon }: URLStatsCardProps) => (
  <div className="bg-white rounded-xl shadow-md p-6 border border-indigo-100 flex items-center">
    <div className="rounded-full bg-indigo-100 p-3 mr-4">
      <Icon className="h-6 w-6 text-indigo-600" />
    </div>
    <div>
      <p className="text-sm text-gray-500">{label}</p>
      <p className="text-2xl font-bold text-gray-900">{count}</p>
    </div>
  </div>
);

export default URLStatsCard;
