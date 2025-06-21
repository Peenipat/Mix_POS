import { useState, useCallback, useEffect } from "react";
import axios from "../lib/axios";
interface Service {
    id: number;
    name: string;
    description: string;
    price: number;
    duration: string;
    Img_path: string;
    Img_name: string
}

export default function ServicePage() {

    const [services, setServices] = useState<Service[]>([]);
    const [loadingServices, setLoadingServices] = useState<boolean>(false);
    const [errorServices, setErrorServices] = useState<string | null>(null);

    const loadServices = useCallback(async () => {
        setLoadingServices(true);
        setErrorServices(null);
        try {
            const res = await axios.get<{ status: string; data: Service[] }>(
                `/barberbooking/tenants/1/branch/1/services`
            );
            if (res.data.status !== "success") {
                throw new Error(res.data.status);
            }
            setServices(res.data.data);
        } catch (err: any) {
            setErrorServices(err.response?.data?.message || err.message || "Failed to load barbers");
        } finally {
            setLoadingServices(false);
        }
    }, []);

    useEffect(() => {
        loadServices();
    }, []);

    return (
        <div className="min-h-screen bg-gray-900 text-gray-200">
            <div className="container mx-auto py-12 px-6">
                <h1 className="text-4xl font-extrabold mb-8">บริการของเรา</h1>

                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
                    {loadingServices && <p>Loading barbers…</p>}
                    {errorServices && <p className="text-red-500">Error loading barbers: {errorServices}</p>}
                    {services.map((service) => {
                        console.log(service)
                        return (
                            <div
                                key={service.id}
                                className="bg-gray-800 rounded-lg shadow-lg hover:shadow-xl transition flex flex-col h-full"
                            >
                                <img
                                    src={`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${service.Img_path}/${service.Img_name}`}
                                    alt={service.name}
                                    className="w-full h-48 object-cover rounded-t-lg"
                                />

                                <div className="p-4 flex flex-col justify-between flex-1">
                                    <div>
                                        <h2 className="text-2xl font-semibold mb-2">{service.name}</h2>
                                        <p className="text-gray-400 text-sm mb-4">{service.description}</p>
                                    </div>
                                    <div className="flex items-center justify-between text-gray-100">
                                        <span className="font-bold text-lg">฿{service.price}</span>
                                        <span className="text-sm bg-gray-700 px-2 py-1 rounded">
                                            {service.duration} นาที
                                        </span>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>

            </div>
        </div>
    );
}
