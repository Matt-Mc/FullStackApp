
export interface BillName {
    en: string;
    fr: string;
  }
  
export interface Bill {
    session: string;
    legisinfo_id: number;
    introduced: string;
    name: BillName; // Changed to lowercase 'n'
    number: string;
    url: string;
  }
  
export interface MP {
    name: string;
    url: string;
    current_party: {
        short_name: {
            en: string;
        }
    };
    current_riding: {
        name: {
            en: string;
        }
        province: string;
    };
    image: string;
  }

