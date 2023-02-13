-- CreateEnum
CREATE TYPE "JackpotNumberType" AS ENUM ('NUMBER', 'POWERBALL');

-- CreateTable
CREATE TABLE "JackpotRow" (
    "id" TEXT NOT NULL,
    "slackWorkspaceId" TEXT NOT NULL,
    "date" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "JackpotRow_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "JackpotNumber" (
    "id" TEXT NOT NULL,
    "slackUserId" TEXT NOT NULL,
    "jackpotRowId" TEXT NOT NULL,
    "numberType" "JackpotNumberType" NOT NULL,
    "number" INTEGER NOT NULL,

    CONSTRAINT "JackpotNumber_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "JackpotNumber_jackpotRowId_numberType_key" ON "JackpotNumber"("jackpotRowId", "numberType");

-- AddForeignKey
ALTER TABLE "JackpotNumber" ADD CONSTRAINT "JackpotNumber_jackpotRowId_fkey" FOREIGN KEY ("jackpotRowId") REFERENCES "JackpotRow"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
