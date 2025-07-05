import Stepper, { Step } from "../../components/Stepper";
import { useState, useEffect } from "react";
export default function AppointmentsPage() {

    const [currentStep, setCurrentStep] = useState(0);
    const isCompleted = currentStep === 4;

    const onStepChange = async (nextStep: number) => {
        setCurrentStep(nextStep);
    };

    return (
        <div className="h-full p-4">
            <main className="flex flex-row gap-4 h-full">
                <div className="border-2 h-full w-1/2 rounded-lg overflow-auto p-2">
                    <h1 className="text-2xl font-bold text-center">จองคิวรายบุคคล</h1>
                    <Stepper
                        step={currentStep}
                        onStepChange={onStepChange}
                        nextButtonText="ถัดไป"
                        backButtonText="ย้อนกลับ"
                        onFinalStepCompleted={() => {
                            setCurrentStep(4);
                        }}
                        className="mt-2 px-6"
                    >
                        <Step>
                            <h2 className="text-start">ขั้นตอนที่ {currentStep + 1}</h2>
                        </Step>
                        <Step>
                            <h2 className="text-start">ขั้นตอนที่ {currentStep + 1}</h2>
                        </Step>
                        <Step>
                            <h2 className="text-start">ขั้นตอนที่ {currentStep + 1}</h2>
                        </Step>
                        <Step>
                            <h2 className="text-start">ขั้นตอนที่ {currentStep + 1}</h2>
                        </Step>
                        {isCompleted && (
                            <div className="text-center py-10">
                                <h2 className="text-2xl font-bold text-green-600">🎉 การจองเสร็จสมบูรณ์</h2>
                                <p className="text-gray-500 mt-2">เราจะติดต่อคุณเร็ว ๆ นี้</p>
                            </div>
                        )}

                    </Stepper>
                </div>
                <div className="border-2 h-full w-1/2 rounded-lg overflow-auto p-2">
                    <h1 className="text-2xl font-bold text-center">จองคิวแบบกลุ่ม</h1>
                    <Stepper
                        step={currentStep}
                        onStepChange={onStepChange}
                        nextButtonText="ถัดไป"
                        backButtonText="ย้อนกลับ"
                        onFinalStepCompleted={() => {
                            setCurrentStep(4);
                        }}
                        className="mt-2 px-6"
                    >
                        <Step>
                            <h2 className="text-start">ขั้นตอนที่ {currentStep + 1}</h2>
                        </Step>
                        <Step>
                            <h2 className="text-start">ขั้นตอนที่ {currentStep + 1}</h2>
                        </Step>
                        <Step>
                            <h2 className="text-start">ขั้นตอนที่ {currentStep + 1}</h2>
                        </Step>
                        <Step>
                            <h2 className="text-start">ขั้นตอนที่ {currentStep + 1}</h2>
                        </Step>
                        {isCompleted && (
                            <div className="text-center py-10">
                                <h2 className="text-2xl font-bold text-green-600">🎉 การจองเสร็จสมบูรณ์</h2>
                                <p className="text-gray-500 mt-2">เราจะติดต่อคุณเร็ว ๆ นี้</p>
                            </div>
                        )}

                    </Stepper>
                </div>
            </main>
        </div>
    );
}
