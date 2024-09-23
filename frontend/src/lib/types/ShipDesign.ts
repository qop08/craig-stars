import type { Cost } from './Cost';
import type { Engine } from './Tech';
import type { TechLevel } from './TechLevel';

export type ShipDesign = {
	id?: number;
	gameId: number;
	createdAt?: Date;
	updatedAt?: Date;
	num?: number;
	playerNum: number;
	originalPlayerNum: number;
	name: string;
	version: number;
	hull: string;
	hullSetNumber: number;
	cannotDelete?: boolean;
	mysteryTrader?: boolean;
	slots: ShipDesignSlot[];
	purpose?: string;
	reportAge?: number;
	spec: Spec;
};

export type ShipDesignSlot = {
	hullComponent: string;
	hullSlotIndex: number;
	quantity: number;
};

export type Bomb = {
	quantity?: number;
	killRate?: number;
	minKillRate?: number;
	structureDestroyRate?: number;
	unterraformRate?: number;
};

export type Spec = {
	armor?: number;
	beamBonus?: number;
	beamDefense?: number;	
	bomber?: boolean;
	bombs?: Bomb[];
	canJump?: boolean;
	canLayMines?: boolean;
	canStealFleetCargo?: boolean;
	canStealPlanetCargo?: boolean;
	cargoCapacity?: number;
	cloakPercent?: number;
	cloakPercentFullCargo?: number;
	cloakUnits?: number;
	colonizer?: boolean;
	cost?: Cost;
	engine: Engine;
	estimatedRange?: number;
	estimatedRangeFull?: number;
	fuelCapacity?: number;
	fuelGeneration?: number;
	hasWeapons?: boolean;
	hullType?: string;
	immuneToOwnDetonation?: boolean;
	initiative?: number;
	mass?: number;
	maxHullMass?: number;
	maxPopulation?: number;
	maxRange?: number;
	mineLayingRateByMineType?: { [mineFieldType: string]: number };
	mineSweep?: number;
	miningRate?: number;
	movement?: number;
	movementBonus?: number;
	movementFull?: number;
	numBuilt?: number;
	numEngines?: number;
	numInstances?: number;
	orbitalConstructionModule?: boolean;
	powerRating?: number;
	reduceCloaking?: number;
	repairBonus?: number;
	retroBombs?: Bomb[];
	safeHullMass?: number;
	safeRange?: number;
	scanner?: boolean;
	scanRange?: number;
	scanRangePen?: number;
	shields?: number;
	smartBombs?: Bomb[];
	spaceDock?: number;
	starbase?: boolean;
	techLevel: TechLevel;
	terraformRate?: number;
	torpedoBonus?: number;
	torpedoJamming?: number;
	weaponSlots?: ShipDesignSlot[];
};
